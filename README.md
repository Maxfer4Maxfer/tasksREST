# tasksREST Service
**TtaskREST Service** is a simple REST API service written on Go.
It can be used as a template or as a project for study purpose.

## Docker-compose Deploy
* [Docker Compose deploy](https://github.com/Maxfer4Maxfer/tasksREST/blob/master/docs/docker-compose-deploy.md)

## Access the service
Create a new tasks and run it for execution. After 2 minutes a status of the created task changed to "finished"
```bash
curl -X POST http://<tasksREST_IP>:8080/task
```

Show information about a particular task
```bash
curl  -X GET http://<tasksREST_IP>:8080/task/{task_id}
```

Show all tasks
```bash
curl  -X GET http://<tasksREST_IP>:8080/tasks
```

## Donations
 If you want to support this project, please consider donating:
 * PayPal: https://paypal.me/MaxFe
