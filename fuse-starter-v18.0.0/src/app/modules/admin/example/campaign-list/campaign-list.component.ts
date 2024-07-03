import { AsyncPipe,CommonModule, CurrencyPipe, NgClass, NgFor, NgIf, NgTemplateOutlet } from '@angular/common';
import { AfterViewInit, ChangeDetectionStrategy, ChangeDetectorRef,NgModule, Component, OnDestroy, OnInit, ViewChild, ViewEncapsulation } from '@angular/core';
import { FormsModule, ReactiveFormsModule, UntypedFormBuilder, UntypedFormControl, UntypedFormGroup, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxChange, MatCheckboxModule } from '@angular/material/checkbox';
import { MatOptionModule, MatRippleModule } from '@angular/material/core';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSelectModule } from '@angular/material/select';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatSort, MatSortModule } from '@angular/material/sort';
import { fuseAnimations } from '@fuse/animations';
import { FuseConfirmationService } from '@fuse/services/confirmation';



import { catchError, debounceTime, map, merge, Observable, of, Subject, switchMap, takeUntil } from 'rxjs';
import { CampaignService } from '../campaign.service';

@Component({
  selector: 'campaignlist',
  standalone: true,
  templateUrl: './campaign-list.component.html',
  styleUrls: ['./campaign-list.component.scss'],
  imports        : [CommonModule,NgIf, MatProgressBarModule, MatFormFieldModule, MatIconModule, MatInputModule, FormsModule, ReactiveFormsModule, MatButtonModule, MatSortModule, NgFor, NgTemplateOutlet, MatPaginatorModule, NgClass, MatSlideToggleModule, MatSelectModule, MatOptionModule, MatCheckboxModule, MatRippleModule, AsyncPipe, CurrencyPipe],

})
export class CampaignListComponent implements OnInit {
createCampaign() {
throw new Error('Method not implemented.');
}

  campaigns: any;

  constructor(private service: CampaignService,) { }

  

ngOnInit(): void {
  this.service.getCampaigns(1, 10).subscribe({
    next: (data) => {
      if (data && data.data && data.data.items) {
        this.campaigns = data.data.items; // Assign mailing list items
        console.log(this.campaigns); // Correctly log the campaigns after they are fetched
      } else {
        console.error('Invalid response structure:', data);
        // Handle unexpected response structure if needed
      }
    },
    error: (error) => {
      console.error('Error fetching mailing lists:', error);
      // Handle error scenario if needed
    }
  });
}

  



deleteProduct() {
throw new Error('Method not implemented.');
}
saveProduct() {
throw new Error('Method not implemented.');
}
  ;
toggleDetails(arg0: number) {
throw new Error('Method not implemented.');
}
selectedProduct: any;
createProduct() {
throw new Error('Method not implemented.');
}


}
