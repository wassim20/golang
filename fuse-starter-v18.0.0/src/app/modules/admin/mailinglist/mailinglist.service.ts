import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { forkJoin, Observable } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { throwError } from 'rxjs';
import { jwtDecode } from 'jwt-decode';
import { AuthService } from 'app/core/auth/auth.service';

@Injectable({
  providedIn: 'root'
})
export class MailinglistService {
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
  
  getMailingLists( page: number = 1, limit: number = 10): Observable<any> {
    const companyID = this.getCompanyID() || '';
    return this.http.get<any>(`${this.baseUrl}/${companyID}/mailinglist`, {
      headers: this.getHeaders(),
      params: {
        page: page.toString(),
        limit: limit.toString()
      }
    });
  }
  getMailingListByID(mailingListId: string): Observable<any> {
    const companyID = this.getCompanyID() || '';
    return this.http.get<any>(`${this.baseUrl}/${companyID}/mailinglist/${mailingListId}`, {
      headers: this.getHeaders()
    });
  }
  
  createMailingList(data: any): Observable<any> {
    const companyID = this.getCompanyID() || '';
    return this.http.post<any>(`${this.baseUrl}/${companyID}/mailinglist`, data, {
      headers: this.getHeaders()
    });
  }
  
  getuser():any{
    const token = this.auth.accessToken;
    const decodedToken: any = jwtDecode(token);
    const userId = decodedToken.userId; // Adjust the key based on your token's payload structure
    return userId;
  }
  
  getlistcreator(userId:string):any{
    console.log("here");
    
    const CompanyID = this.getCompanyID() || '';
    return this.http.get<any>(`http://localhost:8080/api/users/${CompanyID}/${userId}`, {
      headers: this.getHeaders()
    });
    
  }
  
  
  addContactToMailingList(contact: any, mailingListId: string): Observable<any> {
    const companyID = this.getCompanyID() || '';
    const userId = this.getuser();
  
    if (Array.isArray(contact)) {
      console.log(contact);
      
      // If contact is an array (bulk upload), add each contact one by one
      const addContactRequests = contact.map(singleContact =>
        this.http.post<any>(
          `${this.baseUrl}/${companyID}/mailinglist/${mailingListId}/contacts`,
          singleContact,
          { headers: this.getHeaders() }
        )
      );
  
      // Return a combined observable to wait for all requests to complete
      return forkJoin(addContactRequests);
    } else {
      // If contact is a single object, add it directly
      return this.http.post<any>(
        `${this.baseUrl}/${companyID}/mailinglist/${mailingListId}/contacts`,
        contact,
        { headers: this.getHeaders() }
      );
    }
  }
  
  updateMailingList(updatedMailingList: any) {
    const mailinglistupdated = {
    name: updatedMailingList.name,
    description: updatedMailingList.description,
    }
    const companyID = this.getCompanyID() || '';
    return this.http.put<any>(`${this.baseUrl}/${companyID}/mailinglist/${updatedMailingList.id}`, mailinglistupdated, {
      headers: this.getHeaders()
    });
  
  }

  deleteMailingList(mailingListId: string): Observable<any> {
    const companyID = this.getCompanyID() || '';
    return this.http.delete<any>(`${this.baseUrl}/${companyID}/mailinglist/${mailingListId}`, {
      headers: this.getHeaders()
    });
  }
}
