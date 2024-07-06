import { Task } from "src/models/Task";
import { checkStatus } from "./common";

const base_url = process.env.REACT_APP_URL;

export function fetchTasks(): Promise<Task[]> {
  let token = localStorage.getItem("token");
  let url = base_url + "/tasks";

  return fetch(url, {
    headers: { Authorization: `Bearer ${token}` },
  })
    .then(checkStatus)
    .then((response) => {
      if (response) {
        return response.json() as Promise<Task[]>;
      } else {
        return Promise.reject(new Error("Response is undefined"));
      }
    });
}

export function fetchTask(id: string): Promise<Task> {

  let encoded = encodeURIComponent(id);
  const url = base_url + `/tasks/${encoded}`;
  let token = localStorage.getItem("token");

  return fetch(url, { headers: { Authorization: `Bearer ${token}` } })
    .then(checkStatus)
    .then((response) => {
      if (response) {
        return response.json() as Promise<Task>;
      } else {
        return Promise.reject(new Error("Response is undefined"));
      }
    });

}
export function saveNewTask(task: Task): Promise<Task> {
  const url = base_url + `/tasks`;
  const method = "POST";
  return saveTask(url, method, task);
}

export function saveExistingTask(task: Task): Promise<Task> {
  const url = base_url + `/tasks/${encodeURIComponent(task.id)}`;
  const method = "PUT";
  return saveTask(url, method, task);
}
export function saveTask(
  url: string,
  method: string,
  task: Task,
): Promise<Task> {
  let token = localStorage.getItem("token");
  return fetch(url, {
    method: method,
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(task),
  })
    .then(checkStatus)
    .then((response) => {
      if (response) {
        return response.json() as Promise<Task>;
      } else {
        return Promise.reject(new Error("Response is undefined"));
      }
    });
}

export function deleteTask(id: number): Promise<Task | null> {
  let encodedId = encodeURIComponent(id);
  const url = `${base_url}/cards/${encodedId}`;

  let token = localStorage.getItem("token");

  return fetch(url, {
    method: "DELETE",
    headers: { Authorization: `Bearer ${token}` },
  })
    .then(checkStatus)
    .then((response) => {
      if (response) {
        if (response.status === 204) {
          return null;
        }
        return response.json() as Promise<Task>;
      } else {
        return Promise.reject(new Error("Response is undefined"));
      }
    });
}