package triggy

import (
	"errors"
	"fmt"
	"time"

	v1beta1 "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Trigger struct {
	conf Config
}

type Config struct {
	Clientset *kubernetes.Clientset
}

func New(conf Config) (*Trigger, error) {
	if conf.Clientset == nil {
		return nil, errors.New("invalid config")
	}

	t := &Trigger{
		conf: conf,
	}

	return t, nil
}

func (t *Trigger) Run(image string) error {
	err := t.createDeployment(image)
	if err != nil {
		return err
	}
	err = t.scaleDeployment(image)
	if err != nil {
		return err
	}

	return nil
}

func (t *Trigger) createDeployment(image string) error {
	d := deployment(image, 1)
	_, err := t.conf.Clientset.AppsV1beta1().Deployments("default").Create(d)
	if err != nil {
		return err
	}
	fmt.Println("Created new Deployment:", d.Name)
	return nil
}

func (t *Trigger) scaleDeployment(image string) error {
	for i := 2; i < 11; i++ {
		d := deployment(image, i)
		if _, err := t.conf.Clientset.AppsV1beta1().Deployments("default").Update(d); err != nil {
			return err
		}
		fmt.Printf("Deployment scaled by one, to the total of %d instances.\n", i)
		time.Sleep(time.Second * 5)
	}
	fmt.Println("Deployment rollout finished.")
	return nil
}

func deployment(i string, replicas int) *v1beta1.Deployment {
	c := int32(replicas)

	return &v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-123",
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas: &c,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"kubecon": "hello",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "pod-123",
							Image: i,
						},
					},
				},
			},
		},
	}
}
