package oneandone_cloudserver_api

import (
	"github.com/docker/machine/log"
	"net/http"
)

type SharedStorage struct {
	withId
	Size  int    `json:"size"`
	State string `json:"state"`
	withDescription
	CloudPanelId string `json:"cloudpanel_id"`
	SizeUsed     string `json:"size_used"`
	CifsPath     string `json:"cifs_path"`
	NfsPath      string `json:"nfs_path"`
	withName
	CreationDate  string                `json:"creation_date"`
	SharedStorage []SharedStorageServer `json:"servers"`
	withApi
}

type SharedStorageServerPermissions struct {
	SharedStorageServer []SharedStorageServer `json:"servers"`
}

type SharedStorageServer struct {
	Id            string         `json:"id"`
	Name          string         `json:"name"`
	Rights        string         `json:"rights"`
	sharedStorage *SharedStorage `json:"omitempty"`
}

type SharedStorageSettings struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Size        int    `json:"size"`
}

type SharedStorageAccessCredentials struct {
	State               string `json:"state"`
	KerberosContentFile string `json:"kerberos_content_file"`
	UserDomain          string `json:"user_domain"`
	withApi
}

type SharedStorageAccessCredentialsSettings struct {
	Password string `json:"password"`
}

// GET /shared_storages
func (api *API) GetSharedStorages() ([]SharedStorage, error) {
	log.Debug("Requesting information about shared storages")
	result := []SharedStorage{}
	err := api.Client.Get(createUrl(api, SharedStoragesPathSegment), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	for index, _ := range result {
		result[index].api = api
	}
	return result, nil
}

// POST /shared_storages
func (api *API) CreateSharedStorage(configuration SharedStorageSettings) (*SharedStorage, error) {
	log.Debugf("Creating a new shared storage with name '%s'", configuration.Name)
	result := new(SharedStorage)
	err := api.Client.Post(createUrl(api, SharedStoragesPathSegment), configuration, &result, http.StatusAccepted)
	if err != nil {
		return nil, err
	}
	result.api = api
	return result, nil
}

// GET /shared_storages/{id}
func (api *API) GetSharedStorage(sharedStorageId string) (*SharedStorage, error) {
	log.Debugf("Requesting information about the shared storage with the id: '%s'", sharedStorageId)
	result := new(SharedStorage)
	err := api.Client.Get(createUrl(api, SharedStoragesPathSegment, sharedStorageId), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	result.api = api
	return result, nil
}

// DELETE /shared_storages/{id}
func (st *SharedStorage) DeleteSharedStorage() (*SharedStorage, error) {
	log.Debugf("Deleteing shared storage with id: '%s'", st.Id)
	result := new(SharedStorage)
	err := st.api.Client.Delete(createUrl(st.api, SharedStoragesPathSegment, st.Id), &result, http.StatusAccepted)
	if err != nil {
		return nil, err
	}
	result.api = st.api
	return result, nil
}

// PUT /shared_storages/{id}
func (st *SharedStorage) UpdateSharedStorage(configuration SharedStorageSettings) (*SharedStorage, error) {
	log.Debugf("Updateing the shared storage with the id: '%s'", st.Id)
	result := new(SharedStorage)
	err := st.api.Client.Put(createUrl(st.api, SharedStoragesPathSegment, st.Id), configuration, &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	result.api = st.api
	return result, nil
}

// GET /shared_storages/{id}/servers
func (st *SharedStorage) GetSharedStorageServersPermissions() ([]SharedStorageServer, error) {
	log.Debugf("Requesting servers with permissions the the shared storage with the id: '%s'", st.Id)
	result := []SharedStorageServer{}
	err := st.api.Client.Get(createUrl(st.api, SharedStoragesPathSegment, st.Id, "servers"), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	for index, _ := range result {
		result[index].sharedStorage = st
	}
	return result, nil
}

// PUT /shared_storages/{id}/servers
func (st *SharedStorage) UpdateSharedStorageServerPermissions(sharedStorageServerPermissions SharedStorageServerPermissions) (*SharedStorage, error) {
	log.Debugf("Updateing server permissions for the shared storage with the id: '%s'", st.Id)
	result := new(SharedStorage)
	resultError := st.api.Client.Put(createUrl(st.api, SharedStoragesPathSegment, st.Id, "servers"), &sharedStorageServerPermissions, &result, http.StatusAccepted)
	if resultError != nil {
		return nil, resultError
	}
	result.api = st.api
	return result, nil
}

// GET /shared_storages/{id}/servers/{id}
func (st *SharedStorage) GetSharedStorageServersPermission(sharedStorageServerId string) (*SharedStorageServer, error) {
	log.Debugf("Requesting servers permissions for the server: '%s' on the shared storage: '%s' ", sharedStorageServerId, st.Id)
	result := new(SharedStorageServer)
	err := st.api.Client.Get(createUrl(st.api, SharedStoragesPathSegment, st.Id, "servers", sharedStorageServerId), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	result.sharedStorage = st
	return result, nil
}

// DELETE /shared_storages/{id}/servers/{id}
func (sts *SharedStorageServer) DeleteSharedStorageServerPermission() (*SharedStorageServer, error) {
	log.Debugf("Deleting shared storage server permission for the server: '%s' for the shared storage: '%s'", sts.Id, sts.sharedStorage.Id)
	result := new(SharedStorageServer)
	err := sts.sharedStorage.api.Client.Delete(createUrl(sts.sharedStorage.api, SharedStoragesPathSegment, sts.sharedStorage.Id, "servers", sts.Id), &result, http.StatusAccepted)
	if err != nil {
		return nil, err
	}
	result.sharedStorage = sts.sharedStorage
	return result, nil
}

// GET /shared_storages/access
func (api *API) GetSharedStorageAccessCredentials() (*SharedStorageAccessCredentials, error) {
	log.Debugf("Requesting access credentials for the shared storage access")
	result := new(SharedStorageAccessCredentials)
	err := api.Client.Get(createUrl(api, SharedStoragesPathSegment, "access"), &result, http.StatusOK)
	if err != nil {
		return nil, err
	}
	result.api = api
	return result, nil
}

// PUT /shared_storages/access
func (stac *SharedStorageAccessCredentials) UpdateSharedStorageAccessCredentials(sharedStorageAccessCredentialsSettings SharedStorageAccessCredentialsSettings) (*SharedStorageAccessCredentials, error) {
	log.Debugf("Updateing access credentials for the shared storage access")
	result := new(SharedStorageAccessCredentials)
	err := stac.api.Client.Put(createUrl(stac.api, SharedStoragesPathSegment, "access"), &sharedStorageAccessCredentialsSettings, &result, http.StatusAccepted)
	if err != nil {
		return nil, err
	}
	result.api = stac.api
	return result, nil
}
