package git

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/index"
	"github.com/go-git/go-git/v5/plumbing/format/objfile"
	"github.com/go-git/go-git/v5/plumbing/storer"
	gitStorage "github.com/go-git/go-git/v5/storage"
	"go.f110.dev/xerrors"

	"go.f110.dev/mono/go/pkg/collections/set"
	"go.f110.dev/mono/go/pkg/storage"
)

type ObjectStorageInterface interface {
	PutReader(ctx context.Context, name string, data io.Reader) error
	Delete(ctx context.Context, name string) error
	Get(ctx context.Context, name string) (io.ReadCloser, error)
	List(ctx context.Context, prefix string) ([]*storage.Object, error)
}

type ObjectStorageStorer struct {
	backend  ObjectStorageInterface
	rootPath string
}

var _ gitStorage.Storer = &ObjectStorageStorer{}

func NewObjectStorageStorer(b ObjectStorageInterface, rootPath string) *ObjectStorageStorer {
	return &ObjectStorageStorer{backend: b, rootPath: rootPath}
}

func (b *ObjectStorageStorer) Module(name string) (gitStorage.Storer, error) {
	return NewObjectStorageStorer(b.backend, path.Join(b.rootPath, name)), nil
}

func (b *ObjectStorageStorer) Config() (*config.Config, error) {
	file, err := b.backend.Get(context.Background(), path.Join(b.rootPath, "config"))
	if errors.Is(err, storage.ErrObjectNotFound) {
		return config.NewConfig(), nil
	}
	if err != nil {
		return nil, xerrors.WithStack(err)
	}

	conf, err := config.ReadConfig(file)
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	if err := file.Close(); err != nil {
		return nil, xerrors.WithStack(err)
	}
	return conf, nil
}

func (b *ObjectStorageStorer) SetConfig(conf *config.Config) error {
	buf, err := conf.Marshal()
	if err != nil {
		return xerrors.WithStack(err)
	}

	if err := b.backend.PutReader(context.Background(), path.Join(b.rootPath, "config"), bytes.NewReader(buf)); err != nil {
		return xerrors.WithStack(err)
	}
	return nil
}

func (b *ObjectStorageStorer) SetIndex(idx *index.Index) error {
	buf := new(bytes.Buffer)
	if err := index.NewEncoder(buf).Encode(idx); err != nil {
		return xerrors.WithStack(err)
	}

	if err := b.backend.PutReader(context.Background(), path.Join(b.rootPath, "index"), buf); err != nil {
		return xerrors.WithStack(err)
	}
	return nil
}

func (b *ObjectStorageStorer) Index() (*index.Index, error) {
	file, err := b.backend.Get(context.Background(), path.Join(b.rootPath, "index"))
	if err != nil {
		return nil, xerrors.WithStack(err)
	}

	idx := &index.Index{Version: 2}
	if err := index.NewDecoder(file).Decode(idx); err != nil {
		return nil, xerrors.WithStack(err)
	}
	if err := file.Close(); err != nil {
		return nil, xerrors.WithStack(err)
	}
	return idx, nil
}

func (b *ObjectStorageStorer) SetShallow(commits []plumbing.Hash) error {
	buf := new(bytes.Buffer)
	for _, h := range commits {
		if _, err := fmt.Fprintf(buf, "%s\n", h); err != nil {
			return xerrors.WithStack(err)
		}
	}

	if err := b.backend.PutReader(context.Background(), path.Join(b.rootPath, "shallow"), buf); err != nil {
		return xerrors.WithStack(err)
	}
	return nil
}

func (b *ObjectStorageStorer) Shallow() ([]plumbing.Hash, error) {
	file, err := b.backend.Get(context.Background(), path.Join(b.rootPath, "shallow"))
	if err != nil {
		return nil, xerrors.WithStack(err)
	}

	var hash []plumbing.Hash
	s := bufio.NewScanner(file)
	for s.Scan() {
		hash = append(hash, plumbing.NewHash(s.Text()))
	}
	if err := file.Close(); err != nil {
		return nil, xerrors.WithStack(err)
	}
	if err := s.Err(); err != nil {
		return nil, xerrors.WithStack(err)
	}
	return hash, nil
}

func (b *ObjectStorageStorer) SetReference(ref *plumbing.Reference) error {
	buf := new(bytes.Buffer)
	switch ref.Type() {
	case plumbing.SymbolicReference:
		if _, err := fmt.Fprintf(buf, "ref: %s\n", ref.Target()); err != nil {
			return xerrors.WithStack(err)
		}
	case plumbing.HashReference:
		if _, err := fmt.Fprintln(buf, ref.Hash().String()); err != nil {
			return xerrors.WithStack(err)
		}
	}

	if err := b.backend.PutReader(context.Background(), path.Join(b.rootPath, ref.Name().String()), buf); err != nil {
		return xerrors.WithStack(err)
	}
	return nil
}

func (b *ObjectStorageStorer) CheckAndSetReference(new, old *plumbing.Reference) error {
	file, err := b.backend.Get(context.Background(), path.Join(b.rootPath, old.Name().String()))
	if err != nil {
		return xerrors.WithStack(err)
	}
	oldRef, err := b.readReference(file, old.Name().String())
	if err != nil {
		return xerrors.WithStack(err)
	}
	if oldRef.Hash() != old.Hash() {
		return xerrors.New("reference has changed concurrently")
	}

	if err := b.SetReference(new); err != nil {
		return xerrors.WithStack(err)
	}
	return nil
}

func (b *ObjectStorageStorer) Reference(name plumbing.ReferenceName) (*plumbing.Reference, error) {
	file, err := b.backend.Get(context.Background(), path.Join(b.rootPath, name.String()))
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotFound) {
			return nil, plumbing.ErrReferenceNotFound
		}
		return nil, xerrors.WithStack(err)
	}
	ref, err := b.readReference(file, name.String())
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	return ref, nil
}

func (b *ObjectStorageStorer) readReference(f io.ReadCloser, name string) (*plumbing.Reference, error) {
	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	if err := f.Close(); err != nil {
		return nil, xerrors.WithStack(err)
	}
	ref := plumbing.NewReferenceFromStrings(name, strings.TrimSpace(string(buf)))
	return ref, nil
}

func (b *ObjectStorageStorer) IterReferences() (storer.ReferenceIter, error) {
	var refs []*plumbing.Reference
	mark := make(map[plumbing.ReferenceName]struct{})

	// Find refs
	r, err := b.readRefs()
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	for _, v := range r {
		if _, ok := mark[v.Name()]; !ok {
			refs = append(refs, v)
			mark[v.Name()] = struct{}{}
		}
	}

	// Find packed-refs
	packedRefs, err := b.readPackedRefs()
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	for _, v := range packedRefs {
		if _, ok := mark[v.Name()]; !ok {
			refs = append(refs, v)
			mark[v.Name()] = struct{}{}
		}
	}

	// Read HEAD
	ref, err := b.readHEAD()
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	refs = append(refs, ref)

	return storer.NewReferenceSliceIter(refs), nil
}

func (b *ObjectStorageStorer) readRefs() ([]*plumbing.Reference, error) {
	var refs []*plumbing.Reference

	objs, err := b.backend.List(context.Background(), path.Join(b.rootPath, "refs"))
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	for _, v := range objs {
		file, err := b.backend.Get(context.Background(), v.Name)
		if err != nil {
			return nil, xerrors.WithStack(err)
		}
		ref, err := b.readReference(file, strings.TrimPrefix(v.Name, path.Join(b.rootPath, "refs")))
		if err != nil {
			return nil, xerrors.WithStack(err)
		}
		refs = append(refs, ref)
	}

	return refs, nil
}

func (b *ObjectStorageStorer) readPackedRefs() ([]*plumbing.Reference, error) {
	var refs []*plumbing.Reference

	file, err := b.backend.Get(context.Background(), path.Join(b.rootPath, "packed-refs"))
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotFound) {
			return refs, nil
		}
		return nil, xerrors.WithStack(err)
	}
	s := bufio.NewScanner(file)
	for s.Scan() {
		ref, err := b.parsePackedRefsLine(s.Text())
		if err != nil {
			return nil, xerrors.WithStack(err)
		}
		if refs != nil {
			refs = append(refs, ref)
		}
	}
	if err := file.Close(); err != nil {
		return nil, xerrors.WithStack(err)
	}

	return refs, nil
}

func (b *ObjectStorageStorer) parsePackedRefsLine(line string) (*plumbing.Reference, error) {
	switch line[0] {
	case '#', '^':
	default:
		v := strings.Split(line, " ")
		if len(v) != 2 {
			return nil, xerrors.New("git: malformed packed-ref")
		}
		return plumbing.NewReferenceFromStrings(v[1], v[0]), nil
	}

	return nil, nil
}

func (b *ObjectStorageStorer) readHEAD() (*plumbing.Reference, error) {
	file, err := b.backend.Get(context.Background(), path.Join(b.rootPath, "HEAD"))
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	ref, err := b.readReference(file, "HEAD")
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	return ref, nil
}

func (b *ObjectStorageStorer) RemoveReference(name plumbing.ReferenceName) error {
	err := b.backend.Delete(context.Background(), path.Join(b.rootPath, name.String()))
	if err != nil {
		return xerrors.WithStack(err)
	}

	file, err := b.backend.Get(context.Background(), path.Join(b.rootPath, "packed-refs"))
	if err != nil {
		return xerrors.WithStack(err)
	}
	s := bufio.NewScanner(file)
	found := false
	newPackedRefs := new(bytes.Buffer)
	for s.Scan() {
		line := s.Text()
		ref, err := b.parsePackedRefsLine(line)
		if err != nil {
			return xerrors.WithStack(err)
		}
		if ref != nil {
			if ref.Name() == name {
				found = true
				continue
			}
		}
		if _, err := newPackedRefs.WriteString(line); err != nil {
			return xerrors.WithStack(err)
		}
	}
	if err := file.Close(); err != nil {
		return xerrors.WithStack(err)
	}

	if !found {
		// No need to update packed-refs
		return nil
	}

	return b.backend.PutReader(context.Background(), path.Join(b.rootPath, "packed-refs"), newPackedRefs)
}

func (b *ObjectStorageStorer) CountLooseRefs() (int, error) {
	objs, err := b.backend.List(context.Background(), path.Join(b.rootPath, "refs"))
	if err != nil {
		return -1, xerrors.WithStack(err)
	}
	var count int
	mark := make(map[plumbing.ReferenceName]struct{})
	for _, v := range objs {
		file, err := b.backend.Get(context.Background(), v.Name)
		if err != nil {
			return -1, xerrors.WithStack(err)
		}
		ref, err := b.readReference(file, strings.TrimPrefix(v.Name, path.Join(b.rootPath, "refs")))
		if err != nil {
			return -1, xerrors.WithStack(err)
		}
		if _, ok := mark[ref.Name()]; !ok {
			count++
			mark[ref.Name()] = struct{}{}
		}
	}
	return count, nil
}

func (b *ObjectStorageStorer) PackRefs() error {
	r, err := b.readRefs()
	if err != nil {
		return xerrors.WithStack(err)
	}
	if len(r) == 0 {
		return nil
	}
	refSet := set.New()
	for _, v := range r {
		refSet.Add(v)
	}

	packedRefs, err := b.readPackedRefs()
	if err != nil {
		return xerrors.WithStack(err)
	}
	packedRefsSet := set.New()
	for _, v := range packedRefs {
		refSet.Add(v)
		packedRefsSet.Add(v)
	}

	buf := new(bytes.Buffer)
	for _, v := range refSet.ToSlice() {
		ref := v.(*plumbing.Reference)
		if _, err := fmt.Fprintln(buf, ref.String()); err != nil {
			return xerrors.WithStack(err)
		}
	}
	err = b.backend.PutReader(context.Background(), path.Join(b.rootPath, "packed-refs"), buf)
	if err != nil {
		return xerrors.WithStack(err)
	}

	// Delete all loose refs.
	looseRefs := refSet.RightOuter(packedRefsSet)
	for _, v := range looseRefs.ToSlice() {
		ref := v.(*plumbing.Reference)
		err := b.backend.Delete(context.Background(), path.Join(b.rootPath, ref.Name().String()))
		if err != nil {
			return xerrors.WithStack(err)
		}
	}
	return nil
}

func (b *ObjectStorageStorer) NewEncodedObject() plumbing.EncodedObject {
	return &plumbing.MemoryObject{}
}

func (b *ObjectStorageStorer) SetEncodedObject(e plumbing.EncodedObject) (plumbing.Hash, error) {
	switch e.Type() {
	case plumbing.OFSDeltaObject, plumbing.REFDeltaObject:
		return plumbing.ZeroHash, plumbing.ErrInvalidType
	}

	buf := new(bytes.Buffer)
	w := objfile.NewWriter(buf)
	if err := w.WriteHeader(e.Type(), e.Size()); err != nil {
		return plumbing.ZeroHash, xerrors.WithStack(err)
	}
	r, err := e.Reader()
	if err != nil {
		return plumbing.ZeroHash, xerrors.WithStack(err)
	}
	if _, err := io.Copy(w, r); err != nil {
		return plumbing.ZeroHash, xerrors.WithStack(err)
	}
	if err := w.Close(); err != nil {
		return plumbing.ZeroHash, xerrors.WithStack(err)
	}

	hash := w.Hash().String()
	err = b.backend.PutReader(context.Background(), path.Join(b.rootPath, "objects", hash[0:2], hash[2:40]), nil)
	if err != nil {
		return plumbing.ZeroHash, xerrors.WithStack(err)
	}

	return e.Hash(), nil
}

func (b *ObjectStorageStorer) EncodedObject(objectType plumbing.ObjectType, hash plumbing.Hash) (plumbing.EncodedObject, error) {
	// TODO: Read object from pack file
	obj, err := b.getUnpackedEncodedObject(hash)
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	return obj, nil
}

func (b *ObjectStorageStorer) getUnpackedEncodedObject(h plumbing.Hash) (plumbing.EncodedObject, error) {
	file, err := b.backend.Get(context.Background(), path.Join(b.rootPath, "objects", h.String()[0:2], h.String()[2:40]))
	if err != nil {
		return nil, xerrors.WithStack(plumbing.ErrObjectNotFound)
	}

	obj, err := b.readUnpackedEncodedObject(file)
	if err != nil {
		return nil, xerrors.WithStack(err)
	}

	return obj, nil
}

func (b *ObjectStorageStorer) readUnpackedEncodedObject(f io.ReadCloser) (plumbing.EncodedObject, error) {
	obj := b.NewEncodedObject()
	r, err := objfile.NewReader(f)
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	typ, size, err := r.Header()
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	obj.SetType(typ)
	obj.SetSize(size)
	w, err := obj.Writer()
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	if _, err := io.Copy(w, r); err != nil {
		return nil, xerrors.WithStack(err)
	}

	if err := w.Close(); err != nil {
		return nil, xerrors.WithStack(err)
	}
	if err := f.Close(); err != nil {
		return nil, xerrors.WithStack(err)
	}
	return obj, nil
}

func (b *ObjectStorageStorer) IterEncodedObjects(objectType plumbing.ObjectType) (storer.EncodedObjectIter, error) {
	objs, err := b.backend.List(context.Background(), path.Join(b.rootPath, "objects"))
	if err != nil {
		return nil, xerrors.WithStack(err)
	}

	var encodedObjs []plumbing.EncodedObject
	for _, v := range objs {
		s := strings.Split(strings.TrimPrefix(v.Name, b.rootPath), "/")
		if len(s[2]) != 2 || len(s[3]) != 38 {
			continue
		}
		file, err := b.backend.Get(context.Background(), v.Name)
		if err != nil {
			return nil, xerrors.WithStack(err)
		}

		obj, err := b.readUnpackedEncodedObject(file)
		if err != nil {
			return nil, xerrors.WithStack(err)
		}
		encodedObjs = append(encodedObjs, obj)
	}
	return storer.NewEncodedObjectSliceIter(encodedObjs), nil
}

func (b *ObjectStorageStorer) HasEncodedObject(hash plumbing.Hash) error {
	_, err := b.getUnpackedEncodedObject(hash)
	return err
}

func (b *ObjectStorageStorer) EncodedObjectSize(hash plumbing.Hash) (int64, error) {
	obj, err := b.getUnpackedEncodedObject(hash)
	if err != nil {
		return -1, xerrors.WithStack(err)
	}
	return obj.Size(), nil
}