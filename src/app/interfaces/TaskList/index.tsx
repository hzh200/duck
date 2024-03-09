import * as React from 'react';
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

import { Task } from '@/models/Task';
import { columns, ColumnItem } from '@/interfaces/TaskList/TableColumns';

import '@/interfaces/styles/task-list.css';

interface TaskListProps {
  tasks: Array<Task>
}

function TaskList({ tasks }: TaskListProps) {
  
  return (
    <Table>
      <TableHeader>
        <TableRow>
          {columns.map((column, index) => (
            <TableHead key={index} className={column.classes ? column.classes?.join(' ') : ''}>
              {column.title}
            </TableHead>
          ))}
        </TableRow>
      </TableHeader>
      <TableBody className='overflow-scroll'>
        {tasks.map((task, index) => (
          <TableRow key={index}>
            {columns.map((column, index) => (
              <TableCell key={index} className={column.classes ? column.classes?.join(' ') : ''}>
                {(() => {
                  const value = (task as any)[column.accessorKey];
                  if (value !== (null || undefined)) {
                    return column.convertor ? column.convertor(value, task) : value;
                  }
                  return '';
                })()}
              </TableCell>
            ))}
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

export default TaskList;
