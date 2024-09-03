import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { AuthService } from 'app/core/auth/auth.service';
import { Observable } from 'rxjs';
import { jwtDecode } from 'jwt-decode';

@Injectable({
  providedIn: 'root'
})
export class DashboardService {
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
  private baseUrl = 'http://localhost:8080/api';
  //private jwtToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb21wYW55X2lkIjoiYjI3ZWU3N2MtOTA0My00MDAwLWI3ZTAtZjFhOTIwZGEyYzJmIiwiZXhwIjoxNzE5OTY2MDE4LCJpYXQiOjE3MTk3OTMyMTgsInJvbGVzIjpbeyJpZCI6ImMyNmQ0MGE2LTc0MzMtNDg1MC1hYmUyLWIwMzIxYzNkNDA3ZiIsIm5hbWUiOiJNYW5hZ2VyIiwiY29tcGFueV9pZCI6ImIyN2VlNzdjLTkwNDMtNDAwMC1iN2UwLWYxYTkyMGRhMmMyZiJ9XSwidXNlcl9pZCI6IjgzODJlYzZkLTI1OWUtNDQxMy04NjIyLWFiYTFlZWM3MzdlOSJ9.VJn51F28ekCyBpKnwMH-9olqj-FtTp715aBg5LGi2O4'

  constructor(private http: HttpClient,private auth :AuthService) { }

  private getHeaders(): HttpHeaders {
    let headers = new HttpHeaders()
      .set('Authorization', `Bearer ${this.auth.accessToken}`)
      .set('Content-Type', 'application/json');
    return headers;
  }

  getTrackingLogs(campaignID : string, page: number = 1, limit: number = 10): Observable<any> {
    console.log(campaignID);
    const companyID = this.getCompanyID() || '';
  return this.http.get<any>(`${this.baseUrl}/${companyID}/${campaignID}/logs`, {
  headers: this.getHeaders(),
  params: {
    page: page.toString(),
    limit: limit.toString()
  }
  });
  }
  getAllTrackingLogs( page: number = 1, limit: number = 10): Observable<any> {
    const companyID = this.getCompanyID() || '';

    
  return this.http.get<any>(`${this.baseUrl}/${companyID}/logs`, {
  headers: this.getHeaders(),
  params: {
    page: page.toString(),
    limit: limit.toString()
  }
  });
  }


  getCampaigns( page: number = 1, limit: number = 10): Observable<any> {
    const companyID = this.getCompanyID() || '';

    return this.http.get<any>(`${this.baseUrl}/${companyID}/campaigns`, {
      headers: this.getHeaders(),
      params: {
        page: page.toString(),
        limit: limit.toString()
      }
    });
  }
  getCampaign(campaignID : string): Observable<any> {
    console.log("again",campaignID);
    const companyID = this.getCompanyID() || '';
    
    return this.http.get<any>(`${this.baseUrl}/${companyID}/campaigns/${campaignID}`, {
      headers: this.getHeaders()
    });
  }

  getContacts(mailinglistID : string, page: number = 1, limit: number = 10): Observable<any> {
    const companyID = this.getCompanyID() || '';
    return this.http.get<any>(`${this.baseUrl}/${companyID}/mailinglist/${mailinglistID}/contacts`, {
      headers: this.getHeaders(),
      params: {
        page: page.toString(),
        limit: limit.toString()
      }
    });
  }
  getAllContacts( page: number = 1, limit: number = 10): Observable<any> {
    const companyID = this.getCompanyID() || '';
    return this.http.get<any>(`${this.baseUrl}/${companyID}/allcontacts`, {
      headers: this.getHeaders(),
      params: {
        page: page.toString(),
        limit: limit.toString()
      }
    });
  }
  updateChartData(){
    const companyID = this.getCompanyID() || '';
    return this.http.get<any>(`${this.baseUrl}/${companyID}/logs/barchartdata`, {
      headers: this.getHeaders()
    });
  }
  updatePieChartData(){
    const companyID = this.getCompanyID() || '';
    return this.http.get<any>(`${this.baseUrl}/${companyID}/logs/piechartdata`, {
      headers: this.getHeaders()
    });
  }
  updateRadialChartData(){
    const companyID = this.getCompanyID() || '';
    return this.http.get<any>(`${this.baseUrl}/${companyID}/logs/radialchartdata`, {
      headers: this.getHeaders()
    });
  }
  updateLineChartData(){
    const companyID = this.getCompanyID() || '';
    return this.http.get<any>(`${this.baseUrl}/${companyID}/logs/linechartdata`, {
      headers: this.getHeaders()
    });
  }
  updateScatterChartData(){
    const companyID = this.getCompanyID() || '';
    return this.http.get<any>(`${this.baseUrl}/${companyID}/logs/scatterchartdata`, {
      headers: this.getHeaders()
    });
  }
  barChartDataOpens(){
    const companyID = this.getCompanyID() || '';
    return this.http.get<any>(`${this.baseUrl}/${companyID}/logs/barchartopens`, {
      headers: this.getHeaders()
    });
  }
  barChartDataClicks(){
    const companyID = this.getCompanyID() || '';
    return this.http.get<any>(`${this.baseUrl}/${companyID}/logs/barchartclicks`, {
      headers: this.getHeaders()
    });
  }
  bubbleChartDataOpens(){
    const companyID = this.getCompanyID() || '';
    return this.http.get<any>(`${this.baseUrl}/${companyID}/logs/scatterchartopens`, {
      headers: this.getHeaders()
    });
  }
  bubbleChartDataClicks(){
    const companyID = this.getCompanyID() || '';
    return this.http.get<any>(`${this.baseUrl}/${companyID}/logs/scatterchartclicks`, {
      headers: this.getHeaders()
    });
  }

}
