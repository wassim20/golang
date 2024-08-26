import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class DashboardService {

  private companyID ="afa35ff6-4de5-4806-9a21-e0c2453d2834" ;
  private baseUrl = 'http://localhost:8080/api';
  private jwtToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb21wYW55X2lkIjoiYjI3ZWU3N2MtOTA0My00MDAwLWI3ZTAtZjFhOTIwZGEyYzJmIiwiZXhwIjoxNzE5OTY2MDE4LCJpYXQiOjE3MTk3OTMyMTgsInJvbGVzIjpbeyJpZCI6ImMyNmQ0MGE2LTc0MzMtNDg1MC1hYmUyLWIwMzIxYzNkNDA3ZiIsIm5hbWUiOiJNYW5hZ2VyIiwiY29tcGFueV9pZCI6ImIyN2VlNzdjLTkwNDMtNDAwMC1iN2UwLWYxYTkyMGRhMmMyZiJ9XSwidXNlcl9pZCI6IjgzODJlYzZkLTI1OWUtNDQxMy04NjIyLWFiYTFlZWM3MzdlOSJ9.VJn51F28ekCyBpKnwMH-9olqj-FtTp715aBg5LGi2O4'

  constructor(private http: HttpClient) { }

  private getHeaders(): HttpHeaders {
    let headers = new HttpHeaders()
      .set('Authorization', `Bearer ${this.jwtToken}`)
      .set('Content-Type', 'application/json');
    return headers;
  }

  getTrackingLogs(campaignID : string, page: number = 1, limit: number = 10): Observable<any> {
    console.log(campaignID);
    
  return this.http.get<any>(`${this.baseUrl}/${this.companyID}/${campaignID}/logs`, {
  headers: this.getHeaders(),
  params: {
    page: page.toString(),
    limit: limit.toString()
  }
  });
  }
  getAllTrackingLogs( page: number = 1, limit: number = 10): Observable<any> {
    
    
  return this.http.get<any>(`${this.baseUrl}/${this.companyID}/logs`, {
  headers: this.getHeaders(),
  params: {
    page: page.toString(),
    limit: limit.toString()
  }
  });
  }


  getCampaigns( page: number = 1, limit: number = 10): Observable<any> {
    return this.http.get<any>(`${this.baseUrl}/${this.companyID}/campaigns`, {
      headers: this.getHeaders(),
      params: {
        page: page.toString(),
        limit: limit.toString()
      }
    });
  }
  getCampaign(campaignID : string): Observable<any> {
    console.log("again",campaignID);
    
    return this.http.get<any>(`${this.baseUrl}/${this.companyID}/campaigns/${campaignID}`, {
      headers: this.getHeaders()
    });
  }

  getContacts(mailinglistID : string, page: number = 1, limit: number = 10): Observable<any> {
 
    return this.http.get<any>(`${this.baseUrl}/${this.companyID}/mailinglist/${mailinglistID}/contacts`, {
      headers: this.getHeaders(),
      params: {
        page: page.toString(),
        limit: limit.toString()
      }
    });
  }
  getAllContacts( page: number = 1, limit: number = 10): Observable<any> {
 
    return this.http.get<any>(`${this.baseUrl}/${this.companyID}/allcontacts`, {
      headers: this.getHeaders(),
      params: {
        page: page.toString(),
        limit: limit.toString()
      }
    });
  }
  updateChartData(){
    return this.http.get<any>(`${this.baseUrl}/${this.companyID}/logs/barchartdata`, {
      headers: this.getHeaders()
    });
  }
  updatePieChartData(){
    return this.http.get<any>(`${this.baseUrl}/${this.companyID}/logs/piechartdata`, {
      headers: this.getHeaders()
    });
  }
  updateRadialChartData(){
    return this.http.get<any>(`${this.baseUrl}/${this.companyID}/logs/radialchartdata`, {
      headers: this.getHeaders()
    });
  }
  updateLineChartData(){
    return this.http.get<any>(`${this.baseUrl}/${this.companyID}/logs/linechartdata`, {
      headers: this.getHeaders()
    });
  }
  updateScatterChartData(){
    return this.http.get<any>(`${this.baseUrl}/${this.companyID}/logs/scatterchartdata`, {
      headers: this.getHeaders()
    });
  }
  barChartDataOpens(){
    return this.http.get<any>(`${this.baseUrl}/${this.companyID}/logs/barchartopens`, {
      headers: this.getHeaders()
    });
  }
  barChartDataClicks(){
    return this.http.get<any>(`${this.baseUrl}/${this.companyID}/logs/barchartclicks`, {
      headers: this.getHeaders()
    });
  }
  bubbleChartDataOpens(){
    return this.http.get<any>(`${this.baseUrl}/${this.companyID}/logs/scatterchartopens`, {
      headers: this.getHeaders()
    });
  }
  bubbleChartDataClicks(){
    return this.http.get<any>(`${this.baseUrl}/${this.companyID}/logs/scatterchartclicks`, {
      headers: this.getHeaders()
    });
  }

}
