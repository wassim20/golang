/* eslint-disable */
import { FuseNavigationItem } from '@fuse/components/navigation';

export const defaultNavigation: FuseNavigationItem[] = [
    {
        id   : 'dashboard',
        title: 'Dashboard',
        type : 'basic',
        icon : 'heroicons_outline:chart-pie',
        link : '/dashboard'
    },
    {
        id   : 'campaigns',
        title: 'campaigns',
        type : 'basic',
        icon: 'heroicons_outline:chart-bar',
        link : '/campaignlist'
    },
    { // automation
        id: 'automation',
        title: 'Automation',
        type: 'basic',
        icon: 'mat_solid:auto_awesome_motion',
        link: '/automation'
     },
     //mailinglist
        {
            id: 'mailinglist',
            title: 'Mailing List',
            type: 'basic',
            icon: 'mat_solid:email',
            link: '/mailinglist'
        },
        //server
        {
            id: 'server',
            title: 'Server',
            type: 'basic',
            icon: 'heroicons_solid:server-stack',
            link: '/server'
        },
];
export const compactNavigation: FuseNavigationItem[] = [
    {
        id   : 'example',
        title: 'Example',
        type : 'basic',
        icon : 'heroicons_outline:chart-pie',
        link : '/example'
    },
    
];
export const futuristicNavigation: FuseNavigationItem[] = [
    {
        id   : 'example',
        title: 'Example',
        type : 'basic',
        icon : 'heroicons_outline:chart-pie',
        link : '/example'
    },
    
];
export const horizontalNavigation: FuseNavigationItem[] = [
    {
        id   : 'example',
        title: 'Example',
        type : 'basic',
        icon : 'heroicons_outline:chart-pie',
        link : '/example'
    },
    
];
