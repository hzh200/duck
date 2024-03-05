import * as React from 'react';

import { Button } from '@/components/ui/button';
import { MinusIcon, SquareIcon, Cross2Icon } from '@radix-ui/react-icons'

import './css/frame.css';

function Frame() {
  return (
    <div id='frame' className='h-6'>
      <div id='frame-dragger' className='h-6'></div>
      <div id='frame-control' className='align-middle h-6'>
        <Button variant='ghost' size='icon' className='h-6 w-6'>
          <MinusIcon className='h-4 w-4' />
        </Button>
        <Button variant='ghost' size='icon' className='h-6 w-6'>
          <SquareIcon className='h-4 w-4' />
        </Button>
        <Button variant='ghost' size='icon' className='h-6 w-6'>
          <Cross2Icon className='h-4 w-4' />
        </Button>
      </div>
    </div>
  );
}

export default Frame;
