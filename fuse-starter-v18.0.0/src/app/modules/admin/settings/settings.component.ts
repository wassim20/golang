import { CommonModule, CurrencyPipe, NgClass, NgFor, NgIf } from '@angular/common';
import { Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { TranslocoModule } from '@ngneat/transloco';
import { NgApexchartsModule } from 'ng-apexcharts';
import { MatTabsModule } from '@angular/material/tabs';
import { SettingsService } from './settings.service';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { AuthService } from 'app/core/auth/auth.service';
import { UserService } from 'app/core/user/user.service';


@Component({
  selector: 'settings',
  standalone: true,
  imports: [TranslocoModule,MatTabsModule, MatIconModule, MatButtonModule, MatRippleModule, MatMenuModule,
     NgApexchartsModule, NgFor, NgIf,  NgClass, CurrencyPipe,ReactiveFormsModule,CommonModule,FormsModule,
     MatFormFieldModule,MatInputModule,MatSelectModule,MatIconModule],
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss']
})
export class SettingsComponent {
user :any;
users = [];
roles = [];
isFileSelected = false;  // Tracks if a file is selected
selectedFile: File | null = null;  // Stores the selected file
addeduser:any=null;
isAddingRole = false;
profileForm: FormGroup;
assignForm: FormGroup;
addUserForm: FormGroup;
flashMessage: string;


  constructor(private service:SettingsService,private fb: FormBuilder,private _userService:UserService) { }

  ngOnInit(): void {
    this.getUser();
    
    this.addUserForm = this.fb.group({
      userfirstName: ['', Validators.required],
      userlastName: ['', Validators.required],
      useremail: ['', Validators.required],
      userpassword: ['', [Validators.required, Validators.minLength(10), Validators.maxLength(255),]],
      profilePicture: [null],
      
    });
    this.assignForm = this.fb.group({
      user: ['', Validators.required],
      role: ['', Validators.required],
      newRoleName: [''],
    });
    this.loaduserroles();
    this.updateValidators();
  }
  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      this.selectedFile = input.files[0];
      this.isFileSelected = true;
    } else {
      this.selectedFile = null;
      this.isFileSelected = false;
    }
  }
  
  getUser() {
    this.service.getUser().subscribe({
      next: (response) => {
        this.user = response.data;
        // Ensure this happens after user data is fetched
        this.initializeForm();
      },
      error: (error) => {
        console.error('Error fetching user data:', error);
        // Optionally handle the error, e.g., show a notification to the user
      }
    });
  }
  
  initializeForm() {
    this.profileForm = this.fb.group({
      firstname: [this.user?.firstname || '', Validators.required],
      lastname: [this.user?.lastname || '', Validators.required],
      email: [{ value: this.user?.email || '', disabled: true }],
      country: [this.user?.country || '', Validators.required]
    });
  }
  onSubmit(): void {
    if (this.profileForm.valid) {
      const formData = this.profileForm.getRawValue();
  
      // First, handle the form submission for non-file data
      this.service.updateUser(formData).subscribe({
        next: (response) => {
          console.log('User data updated successfully:', response);
  
          // If a file is selected, update the profile picture
          if (this.isFileSelected && this.selectedFile) {
            const pictureFormData = new FormData();
            pictureFormData.append('profilePicture', this.selectedFile);
          
            const entries = (pictureFormData as any).entries();

            for (const [key, value] of entries) {
            console.log(key, value);
            }
          
            this.service.updatePicture(pictureFormData).subscribe({
              next: (pictureResponse) => {
                console.log('Profile picture updated successfully:', pictureResponse);
                
                this._userService.user$.subscribe(user => {
                  user.avatar = pictureResponse.data; // Directly assign the base64 string
                });
                
                this.flashMessage = 'success';
              },
              error: (pictureError) => {
                console.error('Error updating profile picture:', pictureError);
                this.flashMessage = 'error';
              }
            });
            
          } else {
            console.log("No file selected");
            this.flashMessage = 'success';
          }
          
        },
        error: (error) => {
          console.error('Error updating user data:', error);
          this.flashMessage = 'error';
        }
      });
  
      // Reset message after a while
      setTimeout(() => (this.flashMessage = ''), 3000);
    } else {
      // Display error message if form is invalid
      this.flashMessage = 'error';
    }
  }
  
  
  onAddUser(){

    this.service.createUser(this.addUserForm.value).subscribe({
      next: (response) => {
        console.log('User created successfully:', response);
        // Optionally show a success message to the user
        this.flashMessage = 'success';
        // Clear the form
        this.addUserForm.reset();
      },
      error: (error) => {
        console.error('Error creating user:', error);
        // Optionally show an error message to the user
        this.flashMessage = 'error';
      }
    });

    

  }

  loaduserroles(){
    this.service.getRoles().subscribe({
      next: (response) => {
        this.roles = response.data;
        console.log('Roles:', this.roles);
      },
      error: (error) => {
        console.error('Error fetching roles:', error);
        // Optionally handle the error, e.g., show a notification to the user
      }
    });
    this.service.getUsers().subscribe({
      next: (response) => {
        this.users = response.data;
        console.log('Users:', this.users);
      },
      error: (error) => {
        console.error('Error fetching users:', error);
        // Optionally handle the error, e.g., show a notification to the user
      }
    });

  }
  onAssignRole() {
    if (this.assignForm.valid) {
      if (this.isAddingRole) {
        const newRoleName = this.assignForm.get('newRoleName')?.value;
        this.service.createRole(newRoleName).subscribe({
          next: (response) => {
            console.log('Role created successfully:', response);
            this.assignForm.patchValue({ role: response.data.ID });
            this.submitForm();
          },
          error: (error) => {
            console.error('Error creating role:', error);
            this.flashMessage = 'error';
          }
        });
      } else {
        this.submitForm();
      }
    }
  }

  submitForm() {
    const formData = this.assignForm.getRawValue();
    this.service.assignRole(formData).subscribe({
      next: (response) => {
        console.log('Role assigned successfully:', response);
        this.flashMessage = 'success';
      },
      error: (error) => {
        console.error('Error assigning role:', error);
        this.flashMessage = 'error';
      }
    });
  
    setTimeout(() => this.flashMessage = '', 3000);
  }
  
  
updateValidators() {
  const roleControl = this.assignForm.get('role');
  const newRoleNameControl = this.assignForm.get('newRoleName');

  if (this.isAddingRole) {
    newRoleNameControl?.setValidators([Validators.required]);
    roleControl?.clearValidators();
  } else {
    newRoleNameControl?.clearValidators();
    roleControl?.setValidators([Validators.required]);
  }

  newRoleNameControl?.updateValueAndValidity();
  roleControl?.updateValueAndValidity();
}
toggleAddRole() {
  this.isAddingRole = !this.isAddingRole;
  this.updateValidators();
}
 
}
