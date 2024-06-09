import React, { useState } from "react";
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


function TaskPage() {
  const [layout, setLayout] = useState<number[]>([75, 25]);
  const [tasks, setTasks] = useState<Array<Task>>([
    {
      taskNo: 1,
      fileName: "a",
      status: TaskStatus.Successed,
      progress: 1,
      size: 1,
    },
    {
      taskNo: 2,
      fileName: "b",
      status: TaskStatus.Running,
      progress: 66,
      size: 100,
    },
  ]);
  const [choosenFilter, setChoosenFilter] = useState<TaskFilter>(taskFilters[0]);
  const [choosenTaskNos, setChoosenTaskNos] = useState<Array<number>>([]);

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
