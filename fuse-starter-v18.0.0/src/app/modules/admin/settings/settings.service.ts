import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { AuthService } from 'app/core/auth/auth.service';
import { jwtDecode } from 'jwt-decode';

@Injectable({
  providedIn: 'root'
})
export class SettingsService {
  private baseUrl = 'http://localhost:8080/api';
  private getCompanyID(): string | null {
    const token = this.auth.accessToken;
    try {
      const decodedToken: { company_id: string } = jwtDecode(token) as { company_id: string };
      return decodedToken.company_id;
    } catch (error) {
      console.error('Failed to decode JWT token:', error);
      return null;
    }
  }
  
  constructor(private http: HttpClient,private auth:AuthService) { }

  private getHeaders(): HttpHeaders {
    let headers = new HttpHeaders()
    .set('Authorization', `Bearer ${this.auth.accessToken}`)
    .set('Content-Type', 'application/json');
    return headers;
  }

  getuserFromToken():any{
    const token = this.auth.accessToken;
    
    
    const decodedToken: any = jwtDecode(token);
    const userId = decodedToken.user_id; 
    console.log(userId);
    
    return userId;
  }
  
  getUser(){
    const userId = this.getuserFromToken();
    const companyID = this.getCompanyID() || '';
    return this.http.get<any>(`${this.baseUrl}/users/${companyID}/${userId}`, {
      headers: this.getHeaders()
    });
  }

  updateUser(data:any){
    
    const userId = this.getuserFromToken();
    const companyID = this.getCompanyID() || '';
    const userin={
      firstName:data.firstname,
      lastName:data.lastname,
      email:data.email,
      companyID:companyID,
    }
    console.log(userin);
    
    return this.http.put<any>(`${this.baseUrl}/users/${companyID}/${userId}`, userin, {
      headers: this.getHeaders()
    });
  }

  createUser(data:any){
    const companyID = this.getCompanyID() || '';
    const userin={
      firstName:data.userfirstName,
      lastName:data.userlastName,
      email:data.useremail,
      password:data.userpassword,
      companyID:companyID,  
    }
    return this.http.post<any>(`${this.baseUrl}/users/${companyID}`, userin, {
      headers: this.getHeaders()
    });
  }

  getUsers(){
    const companyID = this.getCompanyID() || '';
    return this.http.get<any>(`${this.baseUrl}/users/${companyID}/list`, {
      headers: this.getHeaders()
    });
  }
  getRoles(){
    const companyID = this.getCompanyID() || '';
    return this.http.get<any>(`${this.baseUrl}/roles/${companyID}/list`, {
      headers: this.getHeaders()
    });
  }
  assignRole(data:any){
    const companyID = this.getCompanyID() || '';
    const assign={
      ID:data.user,
      roleID:data.role,
    }
    console.log("heeee",assign);
    
    return this.http.post<any>(`${this.baseUrl}/users/${companyID}/${assign.ID}/roles/${assign.roleID}`,  {
      headers: this.getHeaders()
    });
  }

  createRole(data:any){
    const companyID = this.getCompanyID() || '';
    console.log(data);
    
    const role={
      name:data,
    }
    return this.http.post<any>(`${this.baseUrl}/roles/${companyID}`, role, {
      headers: this.getHeaders()
    });

  }

  updatePicture(formData: FormData) {
    const userId = this.getuserFromToken();
    const companyID = this.getCompanyID() || ''; 
  
    // Create headers with Authorization but without manually setting Content-Type
    const headers = new HttpHeaders({
      'Authorization': `Bearer ${this.auth.accessToken}`
      // Do NOT set 'Content-Type' here; let Angular handle it for FormData
    });
  
    return this.http.put<any>(`${this.baseUrl}/users/${companyID}/${userId}/picture`, formData, { headers });
  }
  
}
