import React from 'react';

import { Separator } from "@/components/ui/separator";
import { Task } from '@/models/Task';

interface TaskInfoProps {
  task: Task | undefined
}

function TaskInfo({ task }: TaskInfoProps) {
  if (!task) {
    return <div></div>
  }
  
  return (
    <div className='h-full'>
      <span>{task.taskName}</span>
      <Separator />
    </div>
  );
}

export default TaskInfo;
