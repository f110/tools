package repoindexer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"go.f110.dev/mono/go/pkg/logger"
	"go.f110.dev/mono/go/pkg/storage"
)

type Manifest struct {
	CreatedAt    time.Time
	Indexes      map[string]string
	ExecutionKey uint64

	filename string
}

func NewManifest(executionKey uint64, indexes map[string]string) Manifest {
	return Manifest{
		CreatedAt:    time.Now(),
		Indexes:      indexes,
		ExecutionKey: executionKey,
		filename:     fmt.Sprintf("manifest_%d.json", executionKey),
	}
}

type ManifestManager struct {
	backend *storage.MinIO
}

func NewManifestManager(backend *storage.MinIO) *ManifestManager {
	return &ManifestManager{backend: backend}
}

func (m *ManifestManager) Update(ctx context.Context, manifest Manifest) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(manifest); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	err := m.backend.Put(ctx, manifest.filename, buf)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	logger.Log.Info("Successfully upload the manifest", zap.String("name", manifest.filename))

	return nil
}

func (m *ManifestManager) GetLatest(ctx context.Context) (Manifest, error) {
	manifest := Manifest{}
	manifests, err := m.backend.List(ctx, "/")
	if err != nil {
		return manifest, xerrors.Errorf(": %w", err)
	}

	latest := int64(0)
	for _, v := range manifests {
		if !strings.HasPrefix(v, "manifest_") {
			continue
		}

		s := strings.TrimSuffix(strings.TrimPrefix(v, "manifest_"), ".json")
		ts, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return manifest, xerrors.Errorf(": %w", err)
		}
		if ts > latest {
			latest = ts
		}
	}
	if latest == 0 {
		return manifest, xerrors.New("repoindexer: Could not found the latest manifest")
	}

	buf, err := m.backend.Get(ctx, fmt.Sprintf("manifest_%d.json", latest))
	if err != nil {
		return manifest, xerrors.Errorf(": %w", err)
	}
	if err := json.NewDecoder(bytes.NewReader(buf)).Decode(&manifest); err != nil {
		return manifest, xerrors.Errorf(": %w", err)
	}
	manifest.filename = fmt.Sprintf("manifest_%d.json", manifest.ExecutionKey)

	return manifest, nil
}

func (m *ManifestManager) Get(ctx context.Context, ts uint64) (Manifest, error) {
	manifest := Manifest{}

	buf, err := m.backend.Get(ctx, fmt.Sprintf("manifest_%d.json", ts))
	if err != nil {
		return manifest, xerrors.Errorf(": %w", err)
	}
	if err := json.NewDecoder(bytes.NewReader(buf)).Decode(&manifest); err != nil {
		return manifest, xerrors.Errorf(": %w", err)
	}
	manifest.filename = fmt.Sprintf("manifest_%d.json", manifest.ExecutionKey)

	return manifest, nil
}

func (m *ManifestManager) FindExpiredManifests(ctx context.Context) ([]Manifest, error) {
	manifests, err := m.backend.List(ctx, "/")
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	timestamps := make([]int64, 0)
	for _, v := range manifests {
		if !strings.HasPrefix(v, "manifest_") {
			continue
		}

		s := strings.TrimSuffix(strings.TrimPrefix(v, "manifest_"), ".json")
		ts, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		timestamps = append(timestamps, ts)
	}
	if len(timestamps) < 3 {
		logger.Log.Debug("Not need cleanup the manifest")
		return nil, nil
	}

	sort.Slice(timestamps, func(i, j int) bool {
		return timestamps[i] > timestamps[j]
	})
	targets := timestamps[2:]
	result := make([]Manifest, 0)
	for _, v := range targets {
		manifest, err := m.Get(ctx, uint64(v))
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		result = append(result, manifest)
	}

	return result, nil
}

func (m *ManifestManager) Delete(ctx context.Context, manifest Manifest) error {
	err := m.backend.Delete(ctx, manifest.filename)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}