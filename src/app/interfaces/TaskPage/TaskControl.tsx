import React from "react";

import {
  PauseIcon,
  PlayIcon,
  TrashIcon,
} from "@radix-ui/react-icons";

import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { TaskFilter, taskFilters } from "@/lib/task-filters";

interface TaskControlProps {
  choosenFilter: TaskFilter
  setChoosenFilter: React.Dispatch<React.SetStateAction<TaskFilter>>
}

function TaskControl({ choosenFilter, setChoosenFilter }: TaskControlProps) {
  return (
    <div className="h-control w-full flex items-center justify-between p-2 border-b">
      <div id="task-control" className="h-full w-fit flex items-center">
        <Button variant="ghost" size="icon" className="h-full w-10">
          <PlayIcon className="h-6 w-6" />
        </Button>
        <Button variant="ghost" size="icon" className="h-full w-10">
          <PauseIcon className="h-6 w-6" />
        </Button>
        <Button variant="ghost" size="icon" className="h-full w-10">
          <TrashIcon className="h-6 w-6" />
        </Button>
      </div>
      <div className="h-full w-fit flex items-center space-x-2">
        <Select onValueChange={(value) => {
          for (const filter of taskFilters) {
            if (filter.name === value) {
              setChoosenFilter(filter)
            }
          }
        }} defaultValue={choosenFilter.name}>
          <SelectTrigger className="w-[120px]">
            <SelectValue placeholder="Select a status" />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectLabel>Status</SelectLabel>
              {taskFilters.map((filter: TaskFilter, index) => (
                <SelectItem key={index} value={filter.name}>{filter.name}</SelectItem>
              ))}
            </SelectGroup>
          </SelectContent>
        </Select>
        <Input className="w-[200px]" type="text" placeholder="name" />
        <Button type="submit">Search</Button>
      </div>
    </div>
  );
}

export default TaskControl;
