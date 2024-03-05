import * as React from 'react';
import Frame from '@/interfaces/Frame';
import MainPanel from '@/interfaces/MainPanel';
import SidePanel from '@/interfaces/SidePanel';

import './css/app.css';

function AppPage() {
  return (
    <div id='app'>
      <Frame />
      <div id='content'>
        <MainPanel />
        <SidePanel />
      </div>
    </div>
  );
}

export default AppPage;
