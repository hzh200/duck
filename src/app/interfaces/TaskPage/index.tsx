import React, { useEffect, useState } from "react";
import TaskList from "@/interfaces/TaskPage/TaskList";
import TaskInfo from "@/interfaces/TaskPage/TaskInfo";
import { Task } from "@/models/Task";
import TaskStatus from "@/models/TaskStatus";
import { taskFilters, TaskFilter } from "@/lib/task-filters";

import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";

import TaskControl from "./TaskControl";

function TaskPage(setting: {[key: string]: any}) {
  const [layout, setLayout] = useState<number[]>([75, 25]);
  const [tasks, setTasks] = useState<Array<Task>>([]);
  const [choosenFilter, setChoosenFilter] = useState<TaskFilter>(taskFilters[0]);
  const [choosenTaskNos, setChoosenTaskNos] = useState<Array<number>>([]);

  useEffect(() => {
    const getTasks = () => {
      fetch(`http://127.0.0.1:${setting["kernelPort"]}/tasks`).then(res => res.json()).then(res => {
        const newTasks: Array<Task> = [];
        for (const task of res["tasks"]) {
          newTasks.push({
            taskNo: task["TaskNo"],
            taskName: task["TaskName"],
            status: task["TaskStatus"],
            progress: task["TaskProgress"],
            size: task["TaskSize"],
          });
        }
        setTasks(newTasks);
      }).catch((_err: Error) => {});
    };

    getTasks(); 
    setInterval(getTasks, 1000);
  }, []);

  return (
    <div className="h-full w-full">
      <ResizablePanelGroup
        direction="horizontal"
        className="min-w-full min-h-full"
        onLayout={(sizes: number[]) => setLayout(sizes)}
      >
        <ResizablePanel defaultSize={layout[0]}>
          <TaskControl choosenFilter={choosenFilter} setChoosenFilter={setChoosenFilter} />
          <div className="h-full p-5">
            <TaskList
              tasks={choosenFilter.filter(tasks)}
              choosen={choosenTaskNos.filter((taskNo) =>
                choosenFilter
                  .filter(tasks)
                  .some((task) => task.taskNo === taskNo)
              )}
              setChoosen={setChoosenTaskNos}
            />
          </div>
        </ResizablePanel>
        <ResizableHandle withHandle  />
        <ResizablePanel defaultSize={layout[1]} minSize={15} maxSize={30}>
          <div className="h-full p-5">
            <TaskInfo
              task={tasks.find(
                (task) =>
                  task.taskNo === choosenTaskNos[choosenTaskNos.length - 1]
              )}
            />
          </div>
        </ResizablePanel>
      </ResizablePanelGroup>
    </div>
  );
}

export default TaskPage;
