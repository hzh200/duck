import React, { useState } from 'react';

import { Task } from '@/models/Task';

import '@/interfaces/styles/task-list.css';
import TaskStatus from '@/models/TaskStatus';
import { TaskFilter, taskFilters } from '../lib/task-filters';
import { cn } from '@/lib/utils';

interface TaskFilterProps {
  choosen: TaskFilter
  setChoosen: React.Dispatch<React.SetStateAction<TaskFilter>>
}

function TaskFilterList({ choosen, setChoosen }: TaskFilterProps) {
  
  return (
    <div>
        {taskFilters.map((filter: TaskFilter, index) => (
          <button key={index} 
            className={cn(
              "flex flex-col items-start gap-2 rounded-lg border p-3 text-left text-sm transition-all hover:bg-accent",
              filter === choosen && "bg-muted"
            )}
            onClick={() => setChoosen(filter)}
          >
            {filter.name}
          </button>
        ))}
    </div>
  );
}

export default TaskFilterList;
