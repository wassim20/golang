import { NO_ERRORS_SCHEMA } from '@angular/compiler';
import { Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { RouterLink } from '@angular/router';
import { FuseVerticalNavigationComponent } from '@fuse/components/navigation';
import { FuseVerticalNavigationAsideItemComponent } from '@fuse/components/navigation/vertical/components/aside/aside.component';
import { FuseVerticalNavigationBasicItemComponent } from '@fuse/components/navigation/vertical/components/basic/basic.component';
import { FuseVerticalNavigationCollapsableItemComponent } from '@fuse/components/navigation/vertical/components/collapsable/collapsable.component';
import { FuseVerticalNavigationDividerItemComponent } from '@fuse/components/navigation/vertical/components/divider/divider.component';
import { FuseVerticalNavigationGroupItemComponent } from '@fuse/components/navigation/vertical/components/group/group.component';
import { FuseVerticalNavigationSpacerItemComponent } from '@fuse/components/navigation/vertical/components/spacer/spacer.component';
import { animate, AnimationBuilder, AnimationPlayer, style } from '@angular/animations';
import { BooleanInput, coerceBooleanProperty } from '@angular/cdk/coercion';
import { ScrollStrategy, ScrollStrategyOptions } from '@angular/cdk/overlay';
import { DOCUMENT, NgFor, NgIf } from '@angular/common';
import { AfterViewInit, ChangeDetectionStrategy, ChangeDetectorRef, ElementRef, EventEmitter, HostBinding, HostListener, Inject, Input, OnChanges, OnDestroy, OnInit, Output, QueryList, Renderer2, SimpleChanges, ViewChild, ViewChildren, ViewEncapsulation } from '@angular/core';
import { NavigationEnd, Router } from '@angular/router';
import { fuseAnimations } from '@fuse/animations';
import { FuseNavigationService } from '@fuse/components/navigation/navigation.service';
import { FuseNavigationItem, FuseVerticalNavigationAppearance, FuseVerticalNavigationMode, FuseVerticalNavigationPosition } from '@fuse/components/navigation/navigation.types';

import { FuseScrollbarDirective } from '@fuse/directives/scrollbar/scrollbar.directive';
import { FuseUtilsService } from '@fuse/services/utils/utils.service';
import { delay, filter, merge, ReplaySubject, Subject, Subscription, takeUntil } from 'rxjs';


@Component({
    selector     : 'landing-home',
    templateUrl  : './home.component.html',
    standalone   : true,
    imports      : [MatButtonModule, RouterLink, MatIconModule,FuseVerticalNavigationComponent,FuseVerticalNavigationAsideItemComponent, FuseVerticalNavigationBasicItemComponent, FuseVerticalNavigationCollapsableItemComponent, FuseVerticalNavigationDividerItemComponent, FuseVerticalNavigationGroupItemComponent, FuseVerticalNavigationSpacerItemComponent,],
    
})
export class LandingHomeComponent
{
    /**
     * Constructor
     */


    constructor()
    {
    }
    navigation = [
        {
            type: 'aside',
          link: '/dashboard',
          label: 'Dashboard',
        },
        {
            type: 'aside',
          link: '/settings',
          label: 'Settings',
        },
      ];
}
