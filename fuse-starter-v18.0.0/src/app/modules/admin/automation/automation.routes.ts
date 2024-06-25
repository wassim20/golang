// automation.routes.ts

import { Routes } from '@angular/router';
import { AutomationComponent } from './automation.component';
import { WorkflowComponent } from './workflow/workflow.component';

const routes: Routes = [
  {
    path: '',
    component: AutomationComponent,
    children: [
      {
        path: 'workflow', 
        component: WorkflowComponent,
      }
    ]
  }
];

export default routes;
