package fdiscovery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Service struct {
	Name     string    `json:"name"`
	Address  string    `json:"address"`
	LastSeen time.Time `json:"last_seen"`
}

type FileSystemDiscovery struct {
	baseDir  string
	mu       sync.RWMutex
	services map[string]*Service
}

func NewFileSystemDiscovery(baseDir string) *FileSystemDiscovery {
	fsd := &FileSystemDiscovery{
		baseDir:  baseDir,
		services: make(map[string]*Service),
	}
	os.MkdirAll(baseDir, 0755)
	go fsd.cleanupLoop()
	return fsd
}

func (f *FileSystemDiscovery) Register(name, address string) error {
	service := &Service{
		Name:     name,
		Address:  address,
		LastSeen: time.Now(),
	}

	data, err := json.Marshal(service)
	if err != nil {
		return err
	}

	filename := filepath.Join(f.baseDir, fmt.Sprintf("%s.json", name))
	return ioutil.WriteFile(filename, data, 0644)
}

func (f *FileSystemDiscovery) Unregister(name string) error {
	filename := filepath.Join(f.baseDir, fmt.Sprintf("%s.json", name))
	return os.Remove(filename)
}

func (f *FileSystemDiscovery) Discover(name string) (*Service, error) {
	filename := filepath.Join(f.baseDir, fmt.Sprintf("%s.json", name))
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var service Service
	err = json.Unmarshal(data, &service)
	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (f *FileSystemDiscovery) Heartbeat(name string) error {
	service, err := f.Discover(name)
	if err != nil {
		return err
	}

	service.LastSeen = time.Now()
	return f.Register(service.Name, service.Address)
}

func (f *FileSystemDiscovery) cleanupLoop() {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		f.cleanup()
	}
}

func (f *FileSystemDiscovery) cleanup() {
	files, err := ioutil.ReadDir(f.baseDir)
	if err != nil {
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filename := filepath.Join(f.baseDir, file.Name())
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			continue
		}

		var service Service
		err = json.Unmarshal(data, &service)
		if err != nil {
			continue
		}

		if time.Since(service.LastSeen) > 1*time.Minute {
			os.Remove(filename)
		}
	}
}
