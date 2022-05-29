# Fullstack To Do App

This is a todo web app with functionalities such as
-	CRUD on To Do items.
-	Sorting and searching for To Do items. 
-	Pagination on To Do item-list.

The project was created with Go on backend and React & TypeScript on frontend.

## Backend:
### Techologies:
**Go** with some **standard libraries (net/http, encoding/json, errors)** and **external packages (gorilla/mux, sirupsen/logrus, google/uuid)**.
 
-	Go is a robust system-level language used for programming across large-scale network servers and big distributed systems. 
-	gorilla/mux used for routing the http requests.
-	logrus used for logging.
-	uuid used for generating unique ids.

### Backend folder hierarchy:
    
      ├── bin               // binary files after compilation
      ├── cmd               // application entry point
      └── pkg               
          ├── api           // api layer
          ├── app           // app
          ├── model         // model corresponds witf frontend
          ├── repository    // interacts with db
          └── version       // informations needed build-time 
      ├── backup.json       // json db
      ├── config.yaml       // configurations
      ├── go.mod            // the root of dependency management
      ├── go.sum            // checksums of the specific module versions
      └── todo.log          // log file
  
### Steps to build and run backend:

After clone this repository, change directory to todo-app-backend then build and execute the project.
For example (remove .exe on linux):
<pre>
cd todo-app-backend
go build -o bin/todo.exe cmd/todo/todo.go
./bin/todo.exe
</pre>

## REST API Documentation

### GET - Get All ToDos By Parameters

Gets all todos by options (paramaters)  
- /api/v1/todos?

Query string parameters:  

| Name  | Data Type | Required/DefaultValue  | Description |
| ------------- | ------------- | ------------- | ------------- |
| page | number | required/1 | Paginated page number of entire list |
| limit | number | required/10 | Number of items that will exist in a single page |
| sortBy | string | required/"dueDate" | Value that will list sorted by |
| sortType | string | required/"desc" | Ascending or descending |
| filter | string | required/"" | Search string |

### GET - Get ToDo By Id

Gets todo by id  
- /api/v1/todos{id}
path variable: id  

### POST - Create To Do

- /api/v1/todos

Request body:  

```
{
  title: string;
  description: string;
  dueDate: string;
}
```

### PUT - Update To Do

- /api/v1/todos

Request body:  

```
{
  id: number;
  title: string;
  description: string;
  dueDate: string;
}
```

### DEL - Delete To Do

- /api/v1/todos{id}

path variable: id

Click [here](https://github.com/yelimot/fullstack-todo-app-frontend) to see the frontend source code.
