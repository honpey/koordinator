package store

import (
	"encoding/json"
)

// MetaManager no need to store info in memory
// as bolt filemap db file
type MetaManager struct {
	db *BoltDB
}

func NewMetaManager() *MetaManager {
	db := NewBoltDB()
	db.Init()
	return &MetaManager{
		db: db,
	}
}

// WritePodSandboxCheckpoint checkpoints the pod level info
func (m *MetaManager) WritePodSandboxCheckpoint(pUID string, pod *PodSandboxCheckpoint) error {
	data, err := json.Marshal(pod)
	if err != nil {
		return err
	}
	m.db.Update(PodBucket, []byte(pUID), data)

	return nil
}

// WriteContainerCheckpoint returns
func (m *MetaManager) WriteContainerCheckpoint(cUID string, container *ContainerCheckpoint) error {
	data, err := json.Marshal(container)
	if err != nil {
		return err
	}
	m.db.Update(ContainerBucket, []byte(cUID), data)
	return nil
}

// GetPodSandboxCheckpoint returns sandbox info
func (m *MetaManager) GetPodSandboxCheckpoint(pUID string) (*PodSandboxCheckpoint, error) {
	pCheckPoint := &PodSandboxCheckpoint{}
	data, err := m.db.Read(PodBucket, []byte(pUID))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, pCheckPoint); err != nil {
		return nil, err
	}
	return pCheckPoint, nil
}

func (m *MetaManager) GetContainerCheckpoint(cUID string) (*ContainerCheckpoint, error) {
	cCheckPoint := &ContainerCheckpoint{}
	data, err := m.db.Read(ContainerBucket, []byte(cUID))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, cCheckPoint); err != nil {
		return nil, err
	}
	return cCheckPoint, nil
}

// DeletePodSandboxCheckpoint delete pod
func (m *MetaManager) DeletePodSandboxCheckpoint(podUID string) error {
	return nil
}

func (m *MetaManager) DeleteContainerCheckpoint(containerUID string) error {
	return nil
}
