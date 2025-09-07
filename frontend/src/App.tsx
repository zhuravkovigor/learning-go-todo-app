import { useMutation, useQuery } from "@tanstack/react-query";
import { useRef, type FC } from "react";

type Todo = {
  id: number;
  title: string;
  completed: boolean;
};

const App: FC = () => {
  const input = useRef<HTMLInputElement>(null);

  const { data, refetch } = useQuery<Todo[]>({
    queryKey: ["todos"],
    queryFn: async () => {
      const response = await fetch("http://localhost:8080/api/todos");
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.json();
    },
  });
  const { mutate } = useMutation({
    mutationFn: async (title: string) => {
      const response = await fetch("http://localhost:8080/api/todo", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ title }),
      });
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.json();
    },

    onSuccess() {
      refetch();
    },
  });
  const { mutate: deleteTodo } = useMutation({
    mutationFn: async (id: number) => {
      await fetch(`http://localhost:8080/api/todo/${id}`, {
        method: "DELETE",
      });
    },

    onSuccess() {
      refetch();
    },
  });
  const { mutate: updateTodo } = useMutation({
    mutationFn: async (todo: Todo) => {
      await fetch(`http://localhost:8080/api/todo/${todo.id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ ...todo, completed: !todo.completed }),
      });
    },

    onSuccess() {
      refetch();
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (input.current) {
      mutate(input.current.value);
      input.current.value = "";
    }
  };

  const handleDeleteTodo = (id: number) => {
    deleteTodo(id);
  };

  const handleDone = (todo: Todo) => {
    updateTodo(todo);
  };

  return (
    <div className="max-w-xl mx-auto py-12">
      <form onSubmit={handleSubmit}>
        <input
          ref={input}
          required
          minLength={3}
          className="bg-zinc-200 w-full p-4 rounded mb-4"
          placeholder="Добавить новый таск"
        />
      </form>

      {data === null && (
        <h2 className="text-2xl text-center py-12 font-bold mb-4">Задач нет</h2>
      )}

      {data?.map((todo) => (
        <div key={todo.id} className="flex gap-4 items-center p-4 border-b">
          <input
            defaultChecked={todo.completed}
            className="bg-red-200"
            type="checkbox"
            onChange={() => handleDone(todo)}
          />
          <h3 className={todo.completed ? "line-through" : ""}>{todo.title}</h3>

          <button
            className="ml-auto bg-red-200 px-3 py-1 rounded cursor-pointer"
            onClick={() => handleDeleteTodo(todo.id)}
          >
            Remove
          </button>
        </div>
      ))}
    </div>
  );
};

export default App;
