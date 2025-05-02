package collector

type ContainerData struct {
	Metrics
	ContainerId      string `json:"container_id"`
	ContainerName    string `json:"container_name"`
	ContainerImage   string `json:"container_image"`
	ContainerStatus  string `json:"container_status"`
	ContainerCreated string `json:"container_created"`
}

func SnapShot() {

}
