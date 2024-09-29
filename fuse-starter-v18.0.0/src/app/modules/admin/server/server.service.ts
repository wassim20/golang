import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { AuthService } from 'app/core/auth/auth.service';
import { jwtDecode } from 'jwt-decode';

@Injectable({
  providedIn: 'root'
})
export class ServerService {
 
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

  getServers(){
    const companyID = this.getCompanyID() || '';
    return this.http.get<any>(`${this.baseUrl}/${companyID}/servers`, {
      headers: this.getHeaders()
    });
  }
  getServerByID(serverId: string){
    const companyID = this.getCompanyID() || '';
    return this.http.get<any>(`${this.baseUrl}/${companyID}/servers/${serverId}`, {
      headers: this.getHeaders()
    });
  }
  updateServer(data: any){
    const companyID = this.getCompanyID() || '';
    return this.http.put<any>(`${this.baseUrl}/${companyID}/servers/${data.id}`, data, {
      headers: this.getHeaders()
    });
  }
  deleteServer(serverId: string){
    const companyID = this.getCompanyID() || '';
    return this.http.delete<any>(`${this.baseUrl}/${companyID}/servers/${serverId}`, {
      headers: this.getHeaders()
    }); 
  }
  createServer(newServer: any) {
    const companyID = this.getCompanyID() || '';
    return this.http.post<any>(`${this.baseUrl}/${companyID}/servers`, newServer, {
      headers: this.getHeaders()
    });
  }

  
  
}
