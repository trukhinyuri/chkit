from functional_tests import chkit
import unittest
import time


class TestDeployment(unittest.TestCase):

    @chkit.test_account
    def test_base(self):
        depl = chkit.Deployment(
            name="functional-test-depl",
            replicas=1,
            containers=[chkit.Container(image="nginx", name="first", limits=chkit.Resources(cpu=10, memory=10))],
        )
        try:
            chkit.create_deployment(depl)
            got_depl = chkit.get_deployment(depl.name)
            self.assertEqual(depl.name, got_depl.name)
            attempts: int
            for i in range(1, 40):
                attempts = i
                pods = chkit.get_pods()
                not_running_pods = [pod for pod in pods if pod.deploy == depl.name and pod.status.phase != "Running"]
                if len(not_running_pods) == 0:
                    break
                time.sleep(15)
            self.assertLessEqual(attempts, 40)
        finally:
            chkit.delete_deployment(name=depl.name)
            time.sleep(5)
            self.assertNotIn(depl.name, [deploy.name for deploy in chkit.get_deployments()])

    @chkit.test_account
    @chkit.with_deployment
    def test_set_image(self, depl: chkit.Deployment):
        chkit.set_image(image="redis", container=depl.containers[0].name, deployment=depl.name)
        got_depl = chkit.get_deployment(depl.name)
        self.assertEqual(got_depl.containers[0].image, "redis")

    @chkit.test_account
    @chkit.with_deployment
    def test_replace_container(self, depl: chkit.Deployment):
        new_container = chkit.Container(
            name=depl.containers[0].name,
            limits=chkit.Resources(cpu=15, memory=15),
            image="redis",
            env={"HELLO": "world"},
        )
        chkit.replace_container(deployment=depl.name, container=new_container)
        got_depl = chkit.get_deployment(depl.name)
        needed_containers = [container for container in got_depl.containers if container.name == new_container.name]
        self.assertGreater(len(needed_containers), 0)
        self.assertEqual(needed_containers[0].env, new_container.env)
        self.assertEqual(needed_containers[0].limits.cpu, new_container.limits.cpu)
        self.assertEqual(needed_containers[0].limits.memory, new_container.limits.memory)
        self.assertEqual(needed_containers[0].image, new_container.image)

    @chkit.test_account
    @chkit.with_deployment
    def test_add_container(self, depl: chkit.Deployment):
        new_container = chkit.Container(
            name="second",
            limits=chkit.Resources(cpu=15, memory=15),
            image="redis",
            env={"HELLO": "world"},
        )
        chkit.add_container(deployment=depl.name, container=new_container)
        got_depl = chkit.get_deployment(depl.name)
        self.assertEqual(len(got_depl.containers), 2)
        needed_containers = [container for container in got_depl.containers if container.name == new_container.name]
        self.assertGreater(len(needed_containers), 0)
        self.assertEqual(needed_containers[0].name, new_container.name)
        self.assertEqual(needed_containers[0].env, new_container.env)
        self.assertEqual(needed_containers[0].limits.cpu, new_container.limits.cpu)
        self.assertEqual(needed_containers[0].limits.memory, new_container.limits.memory)
        self.assertEqual(needed_containers[0].image, new_container.image)

    @chkit.test_account
    @chkit.with_deployment
    @chkit.with_container
    def test_delete_container(self, depl: chkit.Deployment, container: chkit.Container):
        chkit.delete_container(depl.name, container.name)
        got_depl = chkit.get_deployment(depl.name)
        self.assertEqual(len(got_depl.containers), len(depl.containers))
        for i in range(0, len(depl.containers)):
            self.assertEqual(got_depl.containers[i].name, depl.containers[i].name)

    @chkit.test_account
    @chkit.with_deployment
    def test_set_deploy_replicas(self, depl: chkit.Deployment):
        chkit.set_deployment_replicas(deployment=depl.name, replicas=2)
        got_depl = chkit.get_deployment(depl.name)
        self.assertEqual(depl.name, got_depl.name)
        self.assertEqual(got_depl.replicas, 2)

    @chkit.test_account
    @chkit.with_deployment
    @chkit.with_container
    def test_change_deploy_version(self, depl: chkit.Deployment, container: chkit.Container):
        got_depl = chkit.get_deployment(depl.name)
        self.assertIn("1.0.1", got_depl.version)

    @chkit.test_account
    @chkit.with_deployment
    @chkit.with_container
    def test_get_deployment_versions(self, depl: chkit.Deployment, container: chkit.Container):
        deploy_versions = chkit.get_versions(deploy=depl.name)
        self.assertEqual(len(deploy_versions), 2)

    @chkit.test_account
    @chkit.with_deployment
    @chkit.with_container
    def test_run_deployment_version(self, depl: chkit.Deployment, container: chkit.Container):
        chkit.run_version(deploy=depl.name, version="1.0.0")
        time.sleep(5)
        got_depl = chkit.get_deployment(depl.name)
        self.assertEqual(got_depl.version, "1.0.0")
        self.assertEqual(len(got_depl.containers), len(depl.containers))
