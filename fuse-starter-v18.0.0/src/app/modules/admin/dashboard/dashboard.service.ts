import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
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
  private fetchChartData(endpoint: string, startDate?: string, endDate?: string, campaignID?: string) {
    const companyID = this.getCompanyID() || '';
    let params = new HttpParams();

    // Add start and end dates if provided
    if (startDate) {
      params = params.set('startDate', startDate);
    }
    if (endDate) {
      params = params.set('endDate', endDate);
    }
    
    // Add campaignID if provided
    if (campaignID) {
      params = params.set('campaignID', campaignID);
    }

    return this.http.get<any>(`${this.baseUrl}/${companyID}/logs/${endpoint}`, {
      headers: this.getHeaders(),
      params: params
    });
  }

  updateChartData(startDate?: string, endDate?: string, campaignID?: string) {
    return this.fetchChartData('barchartdata', startDate, endDate, campaignID);
  }

  updatePieChartData(startDate?: string, endDate?: string, campaignID?: string) {
    return this.fetchChartData('piechartdata', startDate, endDate, campaignID);
  }

  updateRadialChartData(startDate?: string, endDate?: string, campaignID?: string) {
    return this.fetchChartData('radialchartdata', startDate, endDate, campaignID);
  }

  updateLineChartData(startDate?: string, endDate?: string, campaignID?: string) {
    return this.fetchChartData('linechartdata', startDate, endDate, campaignID);
  }

  updateScatterChartData(startDate?: string, endDate?: string, campaignID?: string) {
    return this.fetchChartData('scatterchartdata', startDate, endDate, campaignID);
  }

  barChartDataOpens(startDate?: string, endDate?: string, campaignID?: string) {
    return this.fetchChartData('barchartopens', startDate, endDate, campaignID);
  }

  barChartDataClicks(startDate?: string, endDate?: string, campaignID?: string) {
    return this.fetchChartData('barchartclicks', startDate, endDate, campaignID);
  }

  bubbleChartDataOpens(startDate?: string, endDate?: string, campaignID?: string) {
    return this.fetchChartData('scatterchartopens', startDate, endDate, campaignID);
  }

  bubbleChartDataClicks(startDate?: string, endDate?: string, campaignID?: string) {
    return this.fetchChartData('scatterchartclicks', startDate, endDate, campaignID);
  }
}


