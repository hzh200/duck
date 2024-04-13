import React from 'react';

import { Separator } from "@/components/ui/separator"

import { cn } from "@/lib/utils";

import '@/interfaces/styles/task-info.css';
import { Task } from '@/models/Task';

interface TaskInfoProps {
  task: Task | undefined
}

function TaskInfo({ task }: TaskInfoProps) {
  if (!task) {
    return <div></div>
  }
  
  return (
    <div>
      <span>{task.fileName}</span>
      <Separator />
    </div>
  );
}

export default TaskInfo;
