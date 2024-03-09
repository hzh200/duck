import React from 'react';

import { Separator } from "@/components/ui/separator"

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";

import '@/interfaces/styles/task-info.css';
import { Task } from '@/models/Task';

interface TaskInfoProps {
  task: Task | null
}

function TaskInfo({ task }: TaskInfoProps) {
  if (task === null) {
    return <div></div>
  }
  
  return (
    <div>
      <span>{task.fileName}</span>
    </div>
  );
}

export default TaskInfo;
