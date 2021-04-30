package consul

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"reflect"

	"github.com/hashicorp/consul/api"
	consulv1alpha1 "go.f110.dev/mono/go/pkg/api/consul/v1alpha1"
	clientset "go.f110.dev/mono/go/pkg/k8s/client/versioned"
	"go.f110.dev/mono/go/pkg/k8s/controllers/controllerutil"
	informers "go.f110.dev/mono/go/pkg/k8s/informers/externalversions"
	"go.f110.dev/mono/go/pkg/storage"
	"golang.org/x/xerrors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	consulv1alpha1listers "go.f110.dev/mono/go/pkg/k8s/listers/consul/v1alpha1"
)

type BackupController struct {
	*controllerutil.ControllerBase

	client            clientset.Interface
	coreClient        kubernetes.Interface
	config            *rest.Config
	runOutsideCluster bool

	backupLister  consulv1alpha1listers.ConsulBackupLister
	serviceLister corev1listers.ServiceLister
	secretLister  corev1listers.SecretLister

	// for testing
	transport http.RoundTripper
}

var _ controllerutil.Controller = &BackupController{}

func NewBackupController(
	coreSharedInformerFactory kubeinformers.SharedInformerFactory,
	sharedInformerFactory informers.SharedInformerFactory,
	coreClient kubernetes.Interface,
	client clientset.Interface,
	config *rest.Config,
	runOutsideCluster bool,
) (*BackupController, error) {
	backupInformer := sharedInformerFactory.Consul().V1alpha1().ConsulBackups()
	serviceInformer := coreSharedInformerFactory.Core().V1().Services()
	secretInformer := coreSharedInformerFactory.Core().V1().Secrets()

	b := &BackupController{
		client:            client,
		coreClient:        coreClient,
		config:            config,
		runOutsideCluster: runOutsideCluster,
		backupLister:      backupInformer.Lister(),
		serviceLister:     serviceInformer.Lister(),
		secretLister:      secretInformer.Lister(),
	}
	b.ControllerBase = controllerutil.NewBase(
		"consul-backup-controller",
		b,
		coreClient,
		[]cache.SharedIndexInformer{backupInformer.Informer()},
		[]cache.SharedIndexInformer{serviceInformer.Informer(), secretInformer.Informer()},
		[]string{},
	)

	return b, nil
}

func (b *BackupController) ObjectToKeys(obj interface{}) []string {
	switch v := obj.(type) {
	case *consulv1alpha1.ConsulBackup:
		key, err := cache.MetaNamespaceKeyFunc(v)
		if err != nil {
			return nil
		}
		return []string{key}
	default:
		return nil
	}
}

func (b *BackupController) GetObject(key string) (runtime.Object, error) {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	backup, err := b.backupLister.ConsulBackups(namespace).Get(name)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return backup, nil
}

func (b *BackupController) UpdateObject(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	backup, ok := obj.(*consulv1alpha1.ConsulBackup)
	if !ok {
		return nil, xerrors.Errorf("unexpected object type: %T", obj)
	}

	updatedBackup, err := b.client.ConsulV1alpha1().ConsulBackups(backup.Namespace).Update(ctx, backup, metav1.UpdateOptions{})
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return updatedBackup, nil
}

func (b *BackupController) Reconcile(ctx context.Context, obj runtime.Object) error {
	backup := obj.(*consulv1alpha1.ConsulBackup)
	updated := backup.DeepCopy()

	consulClient, err := api.NewClient(&api.Config{
		HttpClient: &http.Client{
			Transport: b.transport,
		},
	})
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	now := metav1.Now()
	buf, _, err := consulClient.Snapshot().Save(&api.QueryOptions{})
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	history := &consulv1alpha1.ConsulBackupStatusHistory{
		ExecuteTime: &now,
	}
	if err := b.storeBackupFile(ctx, backup, history, buf, 0, now); err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if history.Succeeded {
		updated.Status.Succeeded = true
		updated.Status.LastSucceededTime = &now
	}
	updated.Status.History = append(updated.Status.History, *history)
	if !reflect.DeepEqual(backup.Status, updated.Status) {
		_, err := b.client.ConsulV1alpha1().ConsulBackups(backup.Namespace).UpdateStatus(ctx, updated, metav1.UpdateOptions{})
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}

	return nil
}

func (b *BackupController) Finalize(ctx context.Context, obj runtime.Object) error {
	panic("implement me")
}

func (b *BackupController) storeBackupFile(
	ctx context.Context,
	backup *consulv1alpha1.ConsulBackup,
	history *consulv1alpha1.ConsulBackupStatusHistory,
	data io.Reader,
	dataSize int64,
	t metav1.Time,
) error {
	switch {
	case backup.Spec.Storage.MinIO != nil:
		spec := backup.Spec.Storage.MinIO

		accessKeySecret, err := b.secretLister.Secrets(backup.Namespace).Get(spec.Credential.AccessKeyID.Name)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
		accessKey, ok := accessKeySecret.Data[spec.Credential.AccessKeyID.Key]
		if !ok {
			return xerrors.Errorf("access key %s not found in %s", spec.Credential.AccessKeyID.Key, accessKeySecret.Name)
		}
		secretAccessKeySecret, err := b.secretLister.Secrets(backup.Namespace).Get(spec.Credential.SecretAccessKey.Name)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
		secretAccessKey, ok := secretAccessKeySecret.Data[spec.Credential.SecretAccessKey.Key]
		if !ok {
			return xerrors.Errorf("secret access key %s not found in %s", spec.Credential.AccessKeyID.Key, accessKeySecret.Name)
		}

		mcOpt := storage.NewMinIOOptions(spec.Service.Name, spec.Service.Namespace, 9000, spec.Bucket, string(accessKey), string(secretAccessKey))
		mcOpt.Transport = b.transport
		mc := storage.NewMinIOStorage(b.coreClient, b.config, mcOpt, b.runOutsideCluster)
		filename := fmt.Sprintf("%s_%d", backup.Name, t.Unix())
		path := spec.Path
		if path[0] == '/' {
			path = path[1:]
		}
		history.Path = filepath.Join(path, filename)
		if err := mc.PutReader(ctx, filepath.Join(path, filename), data, dataSize); err != nil {
			return xerrors.Errorf(": %w", err)
		}

		history.Succeeded = true
		return nil
	case backup.Spec.Storage.GCS != nil:
		spec := backup.Spec.Storage.GCS
		credential, err := b.secretLister.Secrets(backup.Namespace).Get(spec.Credential.ServiceAccountJSONKey.Name)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
		b, ok := credential.Data[spec.Credential.ServiceAccountJSONKey.Key]
		if !ok {
			return xerrors.Errorf("%s is not found in %s", spec.Credential.ServiceAccountJSONKey.Key, spec.Credential.ServiceAccountJSONKey.Name)
		}

		client := storage.NewGCS(b, spec.Bucket)
		filename := fmt.Sprintf("%s_%d", backup.Name, t.Unix())
		if err := client.Put(ctx, data, filepath.Join(spec.Path, filename)); err != nil {
			return xerrors.Errorf(": %w", err)
		}
		history.Path = filepath.Join(spec.Path, filename)

		return nil
	default:
		return xerrors.New("Not configured a storage")
	}
}