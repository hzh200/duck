import * as React from 'react';
import MainPanel from '@/interfaces/MainPanel';
import SidePanel from '@/interfaces/SidePanel';

import './css/globals.css';
import './css/app.css';

function AppPage() {
  return (
    <div id="app">
      <div id="frame">

      </div>
      <MainPanel />
      <SidePanel />
    </div>
  );
}

export default AppPage;
