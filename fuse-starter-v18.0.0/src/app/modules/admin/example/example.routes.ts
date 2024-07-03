import { Routes } from '@angular/router';
import { ExampleComponent } from 'app/modules/admin/example/example.component';
import { CampaignListComponent } from './campaign-list/campaign-list.component';

export default [
    {
        path     : '',
        component: ExampleComponent,
    },
    
] as Routes;
