import { ChangeDetectorRef, Component, Inject } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, FormsModule, ReactiveFormsModule, UntypedFormBuilder, UntypedFormGroup, Validators } from '@angular/forms';
import { ServerService } from './server.service';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MAT_DIALOG_DATA, MatDialog, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { AsyncPipe, CommonModule } from '@angular/common';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatStepperModule } from '@angular/material/stepper';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSelectModule } from '@angular/material/select';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { DateTimePickerModule } from '@syncfusion/ej2-angular-calendars';
import { NgxMatDatetimePickerModule, NgxMatNativeDateModule, NgxMatTimepickerModule } from '@angular-material-components/datetime-picker';
import { MatOptionModule } from '@angular/material/core';

@Component({
  selector: 'app-server',
  standalone: true,
  imports: [
    MatButtonModule, CommonModule, FormsModule, ReactiveFormsModule, MatFormFieldModule,
    MatInputModule, AsyncPipe, CommonModule, MatStepperModule,MatIconModule,
    MatProgressBarModule, MatInputModule, MatSelectModule, MatCheckboxModule, MatDatepickerModule,
    DateTimePickerModule, NgxMatDatetimePickerModule, NgxMatTimepickerModule, NgxMatNativeDateModule,
    MatSnackBarModule,
  ],
  templateUrl: './server.component.html',
  styleUrls: ['./server.component.scss']
})
export class ServerComponent {servers: any[] = [];
  filteredServers = this.servers;
  selectedServerForm: UntypedFormGroup;
  selectedServer: any = null;
  isLoading = false;
  showPassword: boolean = false; // Add this line
  flashMessage: 'success' | 'error' | null = null;
  serverSearchControl = new FormControl('');

  constructor(
    private snackBar: MatSnackBar,
    private service: ServerService,
    private fb: UntypedFormBuilder,
    private _changeDetectorRef: ChangeDetectorRef,
    private dialog: MatDialog
  ) {
    this.serverSearchControl.valueChanges.subscribe(searchText => {
      this.filterServers(searchText);
    });
  }

  ngOnInit(): void {
    // Initialize the form
    
      this.selectedServerForm = this.fb.group({
        id: [''],
        name: ['', Validators.required],
        host: ['', Validators.required],
        port: ['', Validators.required],
        type: ['', Validators.required],
        username: ['', Validators.required],
        password: ['', Validators.required],
      });
    
  

    // Fetch the servers
    this.fetchServers();
  }

  fetchServers(): void {
    this.service.getServers().subscribe({
      next: (data) => {
        if (data && data.data && data.data.items) {
          this.servers = data.data.items;
          this.filteredServers = this.servers;
          console.log('Servers:', this.servers);
        } else {
          console.error('Invalid response structure:', data);
        }
      },
      error: (error) => {
        console.error('Error fetching servers:', error);
      }
    });
  }
  togglePasswordVisibility(): void {
    this.showPassword = !this.showPassword; // Toggle the password visibility
  }

  filterServers(searchText: string): void {
    if (!searchText) {
      // If no search text, reset the filtered servers to the full list
      this.filteredServers = this.servers;
    } else {
      // Otherwise, filter the servers based on the search text
      this.filteredServers = this.servers.filter(server =>
        server.name.toLowerCase().includes(searchText.toLowerCase())
      );
    }
  }

  toggleDetails(serverId: string): void {
    console.log('Toggling details for server:', serverId);

    this.isLoading = true;

    if (this.selectedServer && this.selectedServer.ID === serverId) {
      this.closeDetails();
      return;
    }

    this.service.getServerByID(serverId).subscribe({
      next: (response) => {
        console.log('Server details:', response);
        
        this.selectedServer = response;


        const formValue = {
          id: this.selectedServer.ID, // This should match your server response key
          name: this.selectedServer.Name, // Ensure case matches
          host: this.selectedServer.Host, // Ensure case matches
          port: this.selectedServer.Port, // Ensure case matches
          type: this.selectedServer.Type, // Ensure case matches
          username: this.selectedServer.Username, // Ensure case matches
          password: this.selectedServer.Password, // Ensure case matches
        };
        

        this.selectedServerForm.patchValue(formValue);
        console.log('Server ID:', this.selectedServer.ID);
        this._changeDetectorRef.detectChanges();
        this.isLoading = false;
        console.log('Server form:', this.selectedServer);
        
      },
      error: (error) => {
        console.error('Error fetching server details:', error);
        this.isLoading = false;
      }
    });
  }

  closeDetails(): void {
    this.selectedServer = null;
    this.isLoading = false;
  }

  updateSelectedServer(): void {
    if (this.selectedServerForm.valid) {
      const updatedServer = this.selectedServerForm.value;
      console.log('Updating server:', updatedServer);

      this.service.updateServer(updatedServer).subscribe({
        next: () => {
          this.showFlashMessage('success');
          this.fetchServers();
          this.closeDetails();
        },
        error: (error) => {
          console.error('Error updating server:', error);
        }
      });
    }
  }

  deleteSelectedServer(): void {
    if (this.selectedServer) {
      const serverId = this.selectedServer.ID;
      console.log('Deleting server:', serverId);

      this.service.deleteServer(serverId).subscribe({
        next: () => {
          this.showFlashMessage('success');
          this.fetchServers();
          this.closeDetails();
        },
        error: (error) => {
          console.error('Error deleting server:', error);
        }
      });
    }
  }

  createServer(): void {
    const dialogRef = this.dialog.open(CreateServerDialogComponent, {
      width: '800px',
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.snackBar.open('Server created successfully!', 'Close', {
          duration: 3000,
          verticalPosition: 'top',
        });

        this.fetchServers();
      }
    });
  }

  showFlashMessage(type: 'success' | 'error'): void {
    this.flashMessage = type;
    this._changeDetectorRef.detectChanges();

    setTimeout(() => {
      this.flashMessage = null;
      this._changeDetectorRef.detectChanges();
    }, 3000);
  }}



  @Component({
    selector: 'app-create-server-dialog',
    styleUrls: ['./server.component.scss'],
    templateUrl: './create.component.html',
    standalone: true,
    imports: [
      MatButtonModule, CommonModule, FormsModule, ReactiveFormsModule, MatFormFieldModule,
      MatInputModule, AsyncPipe, CommonModule, MatStepperModule,MatIconModule,
      MatProgressBarModule, MatInputModule, MatSelectModule, MatCheckboxModule, MatDatepickerModule,
      DateTimePickerModule, NgxMatDatetimePickerModule, NgxMatTimepickerModule, NgxMatNativeDateModule,
      MatSnackBarModule,
    ]
  })
  export class CreateServerDialogComponent {
    serverForm: UntypedFormGroup;
  
    constructor(
      private dialogRef: MatDialogRef<CreateServerDialogComponent>,
      private fb: UntypedFormBuilder,
      private snackBar: MatSnackBar,
      private server: ServerService
    ) {
      this.serverForm = this.fb.group({
        name: ['', Validators.required],
        host: ['', Validators.required],
        port: ['', Validators.required],
        type: ['', Validators.required],
        username: ['', Validators.required],
        password: ['', Validators.required]
      });
    }
  
    onSubmit(): void {
      if (this.serverForm.valid) {
        const newServer: any = this.serverForm.value;

        this.server.createServer(newServer).subscribe({
          next: () => {
            this.dialogRef.close(newServer);
            this.snackBar.open('Server created successfully!', 'Close', {
              duration: 3000,
            });
          },
          error: (error) => {
            console.error('Error creating server:', error);
            this.snackBar.open('Failed to create server!', 'Close', {
              duration: 3000,
            });
            this.dialogRef.close();
          }
        });

        
      }
    }
  
    onCancel(): void {
      this.dialogRef.close();
    }
  }
  
