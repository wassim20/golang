import { ChangeDetectorRef, Component, OnInit, ViewChild } from '@angular/core';
import { DashboardService } from '../dashboard.service';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { CommonModule } from '@angular/common'; // Import CommonModule for *ngFor
import { MatOptionModule } from '@angular/material/core'; // Import MatOptionModule for mat-option
import { HttpClientModule } from '@angular/common/http';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import {
  ApexAxisChartSeries,
  ApexChart,
  ApexDataLabels,
  ApexXAxis,
  ApexPlotOptions,
  ApexNonAxisChartSeries,
  ApexResponsive,
  ApexTitleSubtitle,
  ChartComponent,
  NgApexchartsModule,
  ApexStroke,
  ApexYAxis,
  ApexMarkers,
  ApexGrid,
  ApexTooltip,
  ApexFill,
  ApexLegend,
  ApexTheme
} from "ng-apexcharts";
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import html2canvas from 'html2canvas';
import jsPDF from 'jspdf';
//import { ChartOptions } from 'chart.js';

export type BarChartOptions = {
  series: ApexAxisChartSeries;
  chart: ApexChart;
  dataLabels: ApexDataLabels;
  plotOptions: ApexPlotOptions;
  xaxis: ApexXAxis;
  
  title: ApexTitleSubtitle;
};
export type ColumnChartOptions = {
  series: ApexAxisChartSeries;
  chart: ApexChart;
  dataLabels: ApexDataLabels;
  plotOptions: ApexPlotOptions;
  yaxis: ApexYAxis;
  xaxis: ApexXAxis;
  fill: ApexFill;
  title: ApexTitleSubtitle;
};



export type LineChartOptions = {
  series: ApexAxisChartSeries;
  chart: ApexChart;
  xaxis: ApexXAxis;
  stroke: ApexStroke;
  dataLabels: ApexDataLabels;
  markers: ApexMarkers;
  colors: string[];
  yaxis: ApexYAxis;
  grid: ApexGrid;
  legend: ApexLegend;
  title: ApexTitleSubtitle;
};
export type RadialChartOptions = {
  series: ApexNonAxisChartSeries;
  chart: ApexChart;
  labels: string[];
  plotOptions: ApexPlotOptions;
  fill: ApexFill;
  tooltip: ApexTooltip;
};
export type ScatterChartOptions = {
  series: ApexAxisChartSeries;
  chart: ApexChart;
  xaxis: ApexXAxis;
  yaxis: ApexYAxis;
  dataLabels: ApexDataLabels;
  grid: ApexGrid;
  tooltip: ApexTooltip;
  title: ApexTitleSubtitle;
};
export type BubbleChartOptions = {
  series: ApexAxisChartSeries;
  chart: ApexChart;
  xaxis: ApexXAxis;
  yaxis: ApexYAxis;
  plotOptions : ApexPlotOptions;
  fill :ApexFill;
  theme: ApexTheme;
  markers: ApexMarkers;
  dataLabels: ApexDataLabels;
  grid: ApexGrid;
  tooltip: ApexTooltip;
  title: ApexTitleSubtitle;
};

export type PieChartOptions = {
  series: ApexNonAxisChartSeries;
  chart: ApexChart;
  labels: any;
  responsive: ApexResponsive[];
  title: ApexTitleSubtitle;
};


@Component({
  selector: 'campaign-dashboard',
  standalone: true,
  imports: [MatFormFieldModule, MatSelectModule,RouterModule,MatDatepickerModule, FormsModule, ReactiveFormsModule,MatIconModule,
    MatButtonModule,MatTooltipModule,
    CommonModule, MatOptionModule, HttpClientModule,MatProgressSpinnerModule,NgApexchartsModule],
  templateUrl: './campaign-dashboard.component.html',
  styleUrls: ['./campaign-dashboard.component.scss']
})
export class CampaignDashboardComponent {

  public barChartOptions: Partial<BarChartOptions>;
  public barChartOptionsOpen: Partial<ColumnChartOptions>;
  public barChartOptionsClick: Partial<ColumnChartOptions>;
  public pieChartOptions: Partial<PieChartOptions>;
  public lineChartOptions: Partial<LineChartOptions>;
  public radialChartOptions: Partial<RadialChartOptions>;
  public scatterChartOptions: Partial<ScatterChartOptions>;
  public bubbleChartOptionsOpen: Partial<BubbleChartOptions>;
  public bubbleChartOptioncClick: Partial<BubbleChartOptions>;
 
  @ViewChild('barChartOpen') barChartOpen: ChartComponent;
  @ViewChild('barChartClick') barChartClick: ChartComponent;
  @ViewChild('pieChart') pieChart: ChartComponent;
  @ViewChild('radialChart') radialChart: ChartComponent;
  @ViewChild('lineChart') lineChart: ChartComponent;
  @ViewChild('scatterChart') scatterChart: ChartComponent;
  @ViewChild('bubbleChartOpen') bubbleChartOpen: ChartComponent;
  @ViewChild('bubbleChartClick') bubbleChartClick: ChartComponent;

  chartsLoaded: boolean = false;
  barChartInstance: ApexCharts;
  pieChartInstance: ApexCharts;
  lineChartInstance: ApexCharts;
  radialChartInstance: ApexCharts;
  scatterChartInstance: ApexCharts;
  bubbleChartInstance: ApexCharts;
  barChartClickInstance: ApexCharts;
  barChartOpenInstance: ApexCharts;
  bubbleChartOpenInstance: ApexCharts;
  bubbleChartClickInstance: ApexCharts;

  campaigns: any;
  campaignId: any = null;
  contacts: any;
  selectedCampaign: any;
  logs: any;
  reorder : boolean = true;
  stat_type: any = "all";
  dateRangeForm: FormGroup;
  filteredData: any[] = []; // Array to store filtered data
  isLoading: boolean = false;
  selectedStat: string = 'all';
 

   constructor(private fb: FormBuilder,private service : DashboardService,private cdr: ChangeDetectorRef,private router: Router,private route:ActivatedRoute) {  
    this.dateRangeForm = this.fb.group({
      start: [null],
      end: [null]
    });
    const navigation = this.router.getCurrentNavigation();
    const state = navigation?.extras.state as { campaignId: string };
    this.campaignId = state?.campaignId;
 // Initialize bar chart options
 this.barChartOptions = {
  series: [
    {
      name: "Opened",
      data: [0]
    },
    {
      name: "Clicked",
      data: [0]
    },
    {
      name: "Errors",
      data: [0]
    }
  ],
  chart: {
    type: "bar",
    height: 350
  },
  plotOptions: {
    bar: {
      horizontal: true
    }
  },
  dataLabels: {
    enabled: false
  },
  xaxis: {
    categories: ["Opened", "Clicked", "Errors"]
  }
};
this.barChartInstance = new ApexCharts(document.querySelector("#barChart"), this.barChartOptions);


// Initialize pie chart options
this.pieChartOptions = {
  series: [0, 0],
  chart: {
    width: 500,
    type: "pie"
  },
  labels: ["Opened Emails", "Clicked Emails", "Total Contacts"],
  responsive: [{
    breakpoint: 480,
    options: {
      chart: {
        width: 300
      },
      legend: {
        position: "bottom"
      }
    }
  }]
};

this.pieChartInstance = new ApexCharts(document.querySelector("#pieChart"), this.pieChartOptions);

this.lineChartOptions = {
  series: [
    {
      name: "Opened Emails",
      data: []
    },
    {
      name: "Clicked Emails",
      data: []
    }
  ],
  chart: {
    height: 350,
    type: "line",
    dropShadow: {
      enabled: true,
      color: "#000",
      top: 18,
      left: 7,
      blur: 10,
      opacity: 0.2
    },
    toolbar: {
      show: false
    }
  },
  colors: ["#77B6EA", "#545454"],
  dataLabels: {
    enabled: true
  },
  stroke: {
    curve: "smooth"
  },
  title: {
    text: "Email Campaign Tracking",
    align: "left"
  },
  grid: {
    borderColor: "#e7e7e7",
    row: {
      colors: ["#f3f3f3", "transparent"], // takes an array which will be repeated on columns
      opacity: 0.5
    }
  },
  markers: {
    size: 1
  },
  xaxis: {
    categories: [],
    title: {
      text: "Date"
    }
  },
  yaxis: {
    title: {
      text: "Count"
    },
    min: 0,
    max: 10
  },
  legend: {
    position: "top",
    horizontalAlign: "right",
    floating: true,
    offsetY: -25,
    offsetX: -5
  }
};
this.lineChartInstance = new ApexCharts(document.querySelector("#lineChart"), this.lineChartOptions);

this.radialChartOptions = {
  series: [0],
  chart: {
    type: "radialBar",
    offsetY: -20
  },
  plotOptions: {
    radialBar: {
      startAngle: -90,
      endAngle: 90,
      track: {
        background: "#e7e7e7",
        strokeWidth: "97%",
        margin: 5, // margin is in pixels
        dropShadow: {
          enabled: true,
          top: 2,
          left: 0,
          opacity: 0.31,
          blur: 2
        }
      },
      dataLabels: {
        name: {
          show: false
        },
        value: {
          offsetY: -2,
          fontSize: "22px",
          formatter: function(val) {
            return val.toFixed(2) + "%";
          }
        }
      }
    }
  },
  fill: {
    type: "gradient",
    gradient: {
      shade: "light",
      shadeIntensity: 0.4,
      inverseColors: false,
      opacityFrom: 1,
      opacityTo: 1,
      stops: [0, 50, 53, 91]
    }
  },
  labels: ["Filled"],
  
 
};
this.radialChartInstance = new ApexCharts(document.querySelector("#radialChart"), this.radialChartOptions);

this.scatterChartOptions = {
  series: [
    {
      name: "Opened Emails",
      data: []
    },
    // {
    //   name: "Clicked Emails",
    //   data: []
    // }
  ],
  chart: {
    height: 350,
    type: "scatter",
    zoom: {
      type: "xy"
    }
  },
  dataLabels: {
    enabled: false
  },
  grid: {
    xaxis: {
      lines: {
        show: true
      }
    },
    yaxis: {
      lines: {
        show: true
      }
    }
  },
  xaxis: {
    type: "datetime",
    title: {
      text: "Date"
    },
    labels: {
      formatter: function (val) {
        const date = new Date(val);
        return date.toLocaleDateString(); // Format the date as a readable string
      }
    },
    tickAmount: 10,
    max: undefined,
    min: undefined
  },
  yaxis: {
    title: {
      text: "Click Count"
    },
    
  },
  tooltip: {
    enabled: true,
    x: {
      format: 'dd MMM yyyy' // Format the tooltip date
    }
  },
};
this.scatterChartInstance = new ApexCharts(document.querySelector("#scatterChart"), this.scatterChartOptions);
//initialize bar chaart options for open per day of the week
this.barChartOptionsOpen = {
  series: [
    {
      name: "Opened",
      data: []
    },
  ],
  chart: {
    height: 350,
    type: "bar"
  },
  plotOptions: {
    bar: {
      borderRadius :15,
      borderRadiusApplication: "end",
      dataLabels: {
        position: "top" // top, center, bottom
      }
    }
  },
  dataLabels: {
    enabled: true,
    formatter: function(val) {
      return val.toString();
    },
    offsetY: -20,
    style: {
      fontSize: "12px",
      colors: ["#304758"]
    }
  },
  xaxis: {
    categories: [//days of the week
      "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"
      ],
    position: "top",
    labels: {
      offsetY: -18
      
    },
    axisBorder: {
      show: false
    },
    axisTicks: {
      show: false
    },
    
    crosshairs: {
      fill: {
        type: "gradient",
        gradient: {
          colorFrom: "#D8E3F0",
          colorTo: "#BED1E6",
          stops: [0, 100],
          opacityFrom: 0.4,
          opacityTo: 0.5
        }
      }
    },
    tooltip: {
      enabled: true,
      offsetY: -35
    }
  },
  
  yaxis: {
    axisBorder: {
      show: false
    },
    axisTicks: {
      show: false
    },
    labels: {
      show: false,
      formatter: function(val) {
        return val.toString();
      }
    }
  },
  title: {
    text: "Opens per Day",
    floating: false,
    offsetY: 325,
    align: "center",
    style: {
      color: "#444"
    }
  }      
};
this.barChartOpenInstance = new ApexCharts(document.querySelector("#barChartOpen"), this.barChartOptionsOpen);
//initialize bar chaart options for click per day of the week
this.barChartOptionsClick = {
  series: [
    {
      name: "Clicked",
      data: []
    },
  ],
  chart: {
    height: 350,
    type: "bar"
  },
  plotOptions: {
    bar: {
      borderRadius :15,
      borderRadiusApplication: "end",
      dataLabels: {
        position: "top" // top, center, bottom
      }
    }
  },
  dataLabels: {
    enabled: true,
    formatter: function(val) {
      return val.toString();
    },
    offsetY: -20,
    style: {
      fontSize: "12px",
      colors: ["#304758"]
    }
  },
  xaxis: {
    categories: [//days of the week
      "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"
      ],
    position: "top",
    labels: {
      offsetY: -18
    },
    axisBorder: {
      show: false
    },
    axisTicks: {
      show: false
    },
    crosshairs: {
      fill: {
        type: "gradient",
        gradient: {
          colorFrom: "#D8E3F0",
          colorTo: "#BED1E6",
          stops: [0, 100],
          opacityFrom: 0.4,
          opacityTo: 0.5
        }
      }
    },
    tooltip: {
      enabled: true,
      offsetY: -35
    }
  },
  
  yaxis: {
    axisBorder: {
      show: false
    },
    axisTicks: {
      show: false
    },
    labels: {
      show: false,
      formatter: function(val) {
        return val.toString();
      }
    }
  },
  title: {
    text: "clicks per Day",
    floating: false,
    offsetY: 325,
    align: "center",
    style: {
      color: "#444"
    }
  }       
};
this.barChartClickInstance = new ApexCharts(document.querySelector("#barChartClick"), this.barChartOptionsClick);

this.bubbleChartOptionsOpen = {
  series: [{
    name: "Opens",
    data: []
  }],
  chart: {
    height: 350,
    type: "bubble",
    zoom: {
      type: "xy"
    }
  },
  dataLabels: {
    enabled: false
  },
  grid: {
    xaxis: {
      lines: {
        show: true
      }
    },
    yaxis: {
      lines: {
        show: true
      }
    }
  },
  xaxis: {
    categories: ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"],
    title: {
      text: 'Day of Week'
    },
    tickAmount: 6,
    min: 0,
    max: 6,
    labels: {
      formatter: function (val) {
        return ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"][val];
      }
    }
  },
  yaxis: {
    title: {
      text: 'Hour of Day'
    },
    min: 0,
    max: 23,
    tickAmount: 6,
    labels: {
      formatter: function (val) {
        return `${Math.floor(val + 1)}`;
      }
    }
  },
  tooltip: {
    enabled: true,
    x: {
      formatter: (val) => `Day: ${this.getDayName(val)}`
    },
    y: {
      formatter: (val) => `${val} H`
    },
    z: {
      formatter: (val) => `${val} times`
    },
  },
  title: {
    text: 'Email Opens by Time of Day',
    align: 'left'
  },
  plotOptions: {
    bubble: {
      minBubbleRadius: 5,
      maxBubbleRadius: 20,
      zScaling :true
    }
  },
  fill: {
    opacity: 0.8
  },
  theme: {
    mode: 'light', // or 'dark'
    palette: 'palette4',
    monochrome: {
      enabled: false,
      color: '#255aee',
      shadeTo: 'light',
      shadeIntensity: 0.65
    }
  }
};
this.bubbleChartOpenInstance = new ApexCharts(document.querySelector("#bubbleChartOpen"), this.bubbleChartOptionsOpen);

this.bubbleChartOptioncClick = {
  series: [{
    name: "Clicks",
    data: []
  }],chart: {
    height: 350,
    type: "bubble",
    zoom: {
      type: "xy"
    }
  },
  dataLabels: {
    enabled: false
  },
  grid: {
    xaxis: {
      lines: {
        show: true
      }
    },
    yaxis: {
      lines: {
        show: true
      }
    }
  },
  xaxis: {
    categories: ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"],
    title: {
      text: 'Day of Week'
    },
    tickAmount: 6,
    min: 0,
    max: 6,
    labels: {
      formatter: function (val) {
        return ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"][val];
      }
    }
  },
  yaxis: {
    title: {
      text: 'Hour of Day'
    },
    min: 0,
    max: 23,
    tickAmount: 6,
    labels: {
      formatter: function (val) {
        return `${Math.floor(val + 1)}`;
      }
    }
  },
  tooltip: {
    enabled: true,
    x: {
      formatter: (val) => `Day: ${this.getDayName(val)}`
    },
    y: {
      formatter: (val) => `${val} H`
    },
    z: {
      formatter: (val) => `${val} times`
    },
  },
  title: {
    text: 'Email clicks by Time of Day',
    align: 'left'
  },
  plotOptions: {
    bubble: {
      minBubbleRadius: 5,
      maxBubbleRadius: 20,
      zScaling :true
    }
  },
  fill: {
    opacity: 0.8
  },
  theme: {
    mode: 'light', // or 'dark'
    palette: 'palette4',
    monochrome: {
      enabled: false,
      color: '#255aee',
      shadeTo: 'light',
      shadeIntensity: 0.65
    }
  }


};
this.bubbleChartClickInstance = new ApexCharts(document.querySelector("#bubbleChartClick"), this.bubbleChartOptioncClick);

    
    
  }
  
  
  
  
  ngOnInit(): void {
    this.isLoading = true;

    this.route.paramMap.subscribe(params => {
        this.campaignId = params.get('campaignID');
        this.fetchCampaignsAndSelect();
    });
  }
  setDefaultDateRange() {
    const today = new Date();
    const last28Days = new Date(today);
    last28Days.setDate(today.getDate() - 28); // Subtract 28 days
  
    // Log the calculated dates for debugging
    console.log('Today:', today.toISOString());
    console.log('Last 28 Days:', last28Days.toISOString());
  
    // Set the start and end dates in the form as ISO strings
    this.dateRangeForm.patchValue({
        start: last28Days,
        end: today,
    });
  
    // Use the correct order for start and end dates
    const startDateObj = last28Days.toISOString();
    const endDateObj = today.toISOString();
  
    // Log the dates before making API calls for debugging
    console.log('Formatted Start Date:', startDateObj);
    console.log('Formatted End Date:', endDateObj);
  
    // Create an array of promises for chart updates
    const promises = [
        this.updateChartData(startDateObj, endDateObj),
        this.updatePieChartData(startDateObj, endDateObj),
        this.updateLineChartData(startDateObj, endDateObj),
        this.updateRadialChartData(startDateObj, endDateObj),
        this.updateBarChartOpen(startDateObj, endDateObj),
        this.updateBarChartClick(startDateObj, endDateObj),
        this.updateScatterChartData(startDateObj, endDateObj),
        this.updateBubbleChartOpen(startDateObj, endDateObj),
        this.updateBubbleChartClick(startDateObj, endDateObj)
    ];
  
    // Wait for all chart updates to finish
    Promise.all(promises).then(() => {
        console.log('All chart updates are done.');
    }).catch((error) => {
        console.error('Error updating charts:', error);
    });
  }
  
  
  // Method to reset the date range
  resetDateRange() {
    this.setDefaultDateRange(); // Reset to the last 28 days
  }
  filterData(): void {
  
    const startDate = this.dateRangeForm.value.start;
      const endDate = this.dateRangeForm.value.end;
  
      // Format the start and end dates using the formatDate function
      const startDateObj = this.formatDate(startDate);
      const endDateObj = this.formatDate(endDate);
      console.log('Start Date:', startDateObj);
      console.log('End Date:', endDateObj);
      
  
    // Filter logs based on the selected date range using createdAt
    
    
    const promises = [
      this.updateChartData(startDateObj, endDateObj),
      this.updatePieChartData(startDateObj, endDateObj),
      this.updateLineChartData(startDateObj, endDateObj),
      this.updateRadialChartData(startDateObj, endDateObj),
      this.updateBarChartOpen(startDateObj, endDateObj),
      this.updateBarChartClick(startDateObj, endDateObj),
      this.updateScatterChartData(startDateObj, endDateObj),
      this.updateBubbleChartOpen(startDateObj, endDateObj),
      this.updateBubbleChartClick(startDateObj, endDateObj)
  ];
  
  Promise.all(promises).then(() => {
      console.log('All chart updates are done.');
  }).catch((error) => {
      console.error('Error updating charts:', error);
  });
  
    console.log('Filtered Data:', this.filteredData); // Log filtered data or handle it as needed
  }
  formatDate(date: any): string {
    const d = new Date(date);
    const year = d.getFullYear();
    const month = ('0' + (d.getMonth() + 1)).slice(-2);
    const day = ('0' + d.getDate()).slice(-2);
    const hours = ('0' + d.getHours()).slice(-2);
    const minutes = ('0' + d.getMinutes()).slice(-2);
    const seconds = ('0' + d.getSeconds()).slice(-2);
    
    // Adjusting the format to RFC3339
    return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}Z`; // Add 'T' and 'Z' for timezone
  }
private fetchCampaignsAndSelect(): void {
    this.service.getCampaigns(1, 10).subscribe({
        next: (data) => {
            if (data && data.data && data.data.items) {
                this.campaigns = data.data.items;
                this.selectedCampaign = this.campaignId
                    ? this.campaigns.find(campaign => campaign.id === this.campaignId)
                    : 'all';

                if (this.selectedCampaign && this.campaignId) {
                    this.loadCampaignData(this.campaignId);
                } else {
                    this.isLoading = false;
                }
            } else {
                console.error('Invalid response structure:', data);
                this.isLoading = false;
            }
        },
        error: (error) => {
            console.error('Error fetching campaigns:', error);
            this.isLoading = false;
        }
    });
}
private loadCampaignData(campaignId: string): void {
    this.service.getCampaign(campaignId).subscribe({
        next: (data) => {
            if (data && data.data) {
                
                const fetchTrackingDataPromise = this.fetchTrackingData(this.selectedCampaign.id);
                const fetchContactsPromise = this.fetchContacts(this.selectedCampaign.mailingListId);

                Promise.all([fetchTrackingDataPromise, fetchContactsPromise])
                    .then(() => {
                        console.log('Both fetchTrackingData and fetchContacts are done.');
                        this.isLoading = false;
                    })
                    .catch((error) => {
                        console.error('Error in one of the fetch operations:', error);
                        this.isLoading = false;
                    });
            } else {
                console.error('Invalid response structure:', data);
                this.isLoading = false;
            }
        },
        error: (error) => {
            console.error('Error fetching campaign:', error);
            this.isLoading = false;
        }
    });
}
  onCampaignChange(campaign) {
    if (campaign === 'all') {
      this.selectedCampaign = 'all';
      this.router.navigate(['/dashboard']);
    } else {
      this.selectedCampaign = campaign;
       setTimeout(() => {
      this.cdr.detectChanges();
    }, 0);
      this.logs = null;
      this.contacts = null;
      this.isLoading = true;
  
      this.router.navigate(['/dashboard', campaign.id]).then(() => {
        this.selectedCampaign = campaign;
        this.fetchTrackingData(campaign.id);
        this.fetchContacts(campaign.mailingListId);
      });
    }
  }
  fetchContacts(mailingListID: any) {
    this.service.getContacts(mailingListID).subscribe(
        (data) => {
            this.contacts = data?.data?.items?.length > 0 ? data.data.items : null;
            if(this.contacts && this.contacts.length > 0) {
            this.updatePieChartData();
            }else{
              console.log('No contacts to display.');
            }
        },
        (error) => {
            console.error('Error fetching contacts:', error);
            this.contacts = null; // Optionally, also set contacts to null in error scenario
        }
    );
  }
  fetchTrackingData(id: string) {
    this.service.getTrackingLogs(id).subscribe(
        (data) => {
            this.logs = data?.data?.items?.length > 0 ? data.data.items : null;
            // if logs are not empty 
            if (this.logs && this.logs.length > 0) {
                const promises = [
                    this.updateChartData(),
                    this.updateLineChartData(),
                    this.updateRadialChartData(),
                    this.updateBarChartOpen(),
                    this.updateBarChartClick(),
                    this.updateScatterChartData(),
                    this.updateBubbleChartOpen(),
                    this.updateBubbleChartClick()
                ];

                Promise.all(promises).then(() => {
                    console.log('All chart updates are done.');
                    this.isLoading = false; // Set isLoading to false after all chart updates are done
                }).catch((error) => {
                    console.error('Error updating charts:', error);
                });
            } else {
              this.isLoading = false;
                console.log('No logs to display.');
            }
        },
        (error) => {
            console.error('Error fetching logs:', error);
            this.logs = null; // Optionally, also set logs to null in error scenario
        }
    );
}
updateChartData(startDate?: string, endDate?: string): void {
  let totalOpened = 0;
  let totalClicked = 0;
  let totalErrors = 0;
  
  this.service.updateChartData(startDate,endDate,this.campaignId).subscribe((data) => {
    
    totalOpened = data.data.opened;
    totalClicked = data.data.clicked;
    totalErrors = data.data.error;
    
    // Adjusting the structure to match a single dataset for all categories
    this.barChartOptions.series = [
      {
        name: "",
        data: [totalOpened, totalClicked, totalErrors]
      },
    ];
  });
   setTimeout(() => {
       setTimeout(() => {
      this.cdr.detectChanges();
    }, 0);
    }, 0);
  
  
}
updateBubbleChartClick(startDate?: string, endDate?: string): void {
  this.service.bubbleChartDataClicks(startDate,endDate,this.campaignId).subscribe((response) => {
    

    // Update the bubble chart series data using the data from the backend
    this.bubbleChartOptioncClick.series = [{
      name: 'Clicks',
      data: response.data.map(item => ({
        x: item.x,  // Day of the week (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
        y: item.y,  // Hour of the day (0 = midnight, 1 = 1 AM, ..., 23 = 11 PM)
        z: item.z   // Bubble size (1 = small, 2 = medium, 3 = large)
      }))
    }];

    // Log the series data for debugging purposes
    console.log(this.bubbleChartOptioncClick.series[0].data);
     setTimeout(() => {
       setTimeout(() => {
      this.cdr.detectChanges();
    }, 0);
    }, 0);

    // Trigger change detection to update the chart
  },
  (error) => {
    console.error('Error fetching bubble chart click data:', error);
    // Handle the error appropriately, e.g., display an error message
  });
}

updateBubbleChartOpen(startDate?: string, endDate?: string): void {
  this.service.bubbleChartDataOpens(startDate,endDate,this.campaignId).subscribe((response) => {
    

    // Update the bubble chart series data using the data from the backend
    this.bubbleChartOptionsOpen.series = [{
      name: 'Opens',
      data: response.data.map(item => ({
        x: item.x,  // Day of the week (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
        y: item.y,  // Hour of the day (0 = midnight, 1 = 1 AM, ..., 23 = 11 PM)
        z: item.z   // Bubble size (1 = small, 2 = medium, 3 = large)
      }))
    }];
    // Trigger change detection to update the chart
     setTimeout(() => {
       setTimeout(() => {
      this.cdr.detectChanges();
    }, 0);
    }, 0);
  },
  (error) => {
    console.error('Error fetching bubble chart data:', error);
    // Handle the error appropriately, e.g., display an error message
  });
   setTimeout(() => {
      this.cdr.detectChanges();
    }, 0);
}

updateBarChartClick(startDate?: string, endDate?: string) {
  this.service.barChartDataClicks(startDate, endDate, this.campaignId).subscribe((data) => {
    console.log('Click Data:', data); // Log the click data

    // Check if the chart is initialized
    console.log('Bar Chart Click Instance:', this.barChartClick);

    if (data && data.data && this.barChartClick) {
      this.barChartClick.updateSeries([{
        name: 'Clicks',
        data: data.data
      }], true);
      this.cdr.detectChanges();
    }
  }, error => {
    console.error("Error fetching click data", error);
  });
}

updateBarChartOpen(startDate?: string, endDate?: string) {
  this.service.barChartDataOpens(startDate, endDate, this.campaignId).subscribe((data) => {
    console.log('Open Data:', data); // Log the open data

    // Check if the chart is initialized
    console.log('Bar Chart Open Instance:', this.barChartOpen);

    if (data && data.data && this.barChartOpen) {
      this.barChartOpen.updateSeries([{
        name: 'Opens',
        data: data.data
      }], true);
      this.cdr.detectChanges();
    }
  }, error => {
    console.error("Error fetching open data", error);
  });
}



updateScatterChartData(startDate?: string, endDate?: string): void {
this.service.updateScatterChartData(startDate, endDate,this.campaignId).subscribe((data) => {
// Check if openedData and clickedData are available
const openedData = data.data.openedData || [];
const clickedData = data.data.clickedData || [];

// Update the scatter chart series data using the data from the backend
this.scatterChartOptions.series = [
{
  name: "Opened Emails",
  data: openedData.map(item => ({
    x: item.x,  // Timestamp already provided by the backend
    y: item.y,  // Click count provided by the backend
    recipientEmail: item.recipientEmail || 'Unknown'
  }))
},
{
  name: "Clicked Emails",
  data: clickedData.map(item => ({
    x: item.x,  // Timestamp already provided by the backend
    y: item.y,  // Click count provided by the backend
    recipientEmail: item.recipientEmail || 'Unknown'
  }))
}
];

// Update tooltip to use the custom data
this.scatterChartOptions.tooltip = {
enabled: true,
custom: ({ seriesIndex, dataPointIndex, w }) => {
  const recipientEmail = w.config.series[seriesIndex].data[dataPointIndex]?.recipientEmail || 'Unknown';
  const clickCount = w.config.series[seriesIndex].data[dataPointIndex]?.y || 0;
  return `<div class="apexcharts-tooltip-title">Recipient Email: ${recipientEmail}</div>
  <div class="apexcharts-tooltip-title">Click Count: ${clickCount}</div>`;
}
};

// Trigger change detection to update the chart
 setTimeout(() => {
       setTimeout(() => {
      this.cdr.detectChanges();
    }, 0);
    }, 0);
},
(error) => {
console.error('Error fetching scatter chart data:', error);
// Optionally handle the error, such as displaying a notification to the user
});
}
updateRadialChartData(startDate?: string, endDate?: string) {
const totalLogs = this.logs.length;

this.service.updateRadialChartData(startDate,endDate,this.campaignId).subscribe((data) => {



this.radialChartOptions.tooltip = {
  enabled: true,
  y: {
    formatter: () => {
      
      
      return `Opened Emails: ${data.data.openedLogs}`;
    },
  },
  
};
// Update the radial chart with the calculated percentage
this.radialChartOptions.series = [data.data.openedPercentage];
 setTimeout(() => {
       setTimeout(() => {
      this.cdr.detectChanges();
    }, 0);
    }, 0);
});
 setTimeout(() => {
       setTimeout(() => {
      this.cdr.detectChanges();
    }, 0);
    }, 0);
}
updateLineChartData(startDate?: string, endDate?: string): void {
this.service.updateLineChartData(startDate, endDate,this.campaignId).subscribe((data) => {
  if (this.lineChartOptions?.chart) {
    const allDates = data.data.allDates ? data.data.allDates.map(date => `${date}`) : [];
    const totals = data.data.totals || 0;

    this.lineChartOptions = {
      ...this.lineChartOptions,
      xaxis: {
        categories: allDates,
      },
      series: [
        { name: "Opened Emails", data: data.data.openedSeriesData || [] },
        { name: "Clicked Emails", data: data.data.clickedSeriesData || [] }
      ],
      title: {
        text: `Email Campaign Tracking (Total Sent: ${totals})`,
        align: "left"
      }
    };

     setTimeout(() => {
       setTimeout(() => {
      this.cdr.detectChanges();
    }, 0);
    }, 0);  // Only detect changes once at the end
  } else {
    console.error("Line chart is not initialized");
  }
}, (error) => {
  console.error('Error fetching line chart data:', error);
});
}
updatePieChartData(startDate?: string, endDate?: string): void {


const totalContacts = this.contacts.length;

this.service.updatePieChartData(startDate,endDate,this.campaignId).subscribe(
  (data) => {
    
    const openedEmails = data.data.opened_count || 0;
    const clickedEmails = data.data.clicked_count || 0;

    this.pieChartOptions = {
      ...this.pieChartOptions,
      labels: ["Opened Emails", "Clicked Emails", "Total Contacts: " + totalContacts],
      series: [openedEmails, clickedEmails],
    };
    
     setTimeout(() => {
      this.cdr.detectChanges();
    }, 0);  // Only detect changes once at the end
  },
  (error) => {
    console.error('Error fetching data from backend:', error);
  }
);
}
getDayName(dayIndex) {
  const days = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
  return days[dayIndex];
}
//stat change 
statchange(type: string) {
  this.selectedStat = type;
  this.reorder = false; // Temporarily hide the charts
  this.cdr.detectChanges(); // Trigger change detection

  // Delay to ensure the DOM updates
  setTimeout(() => {
    // Set the state type based on the selected type
    if (type === 'open') {
      this.stat_type = 'open';
      this.updateBarChartOpen(); // Call function to update bar chart for opens
      this.updateBubbleChartOpen(); // Call function to update bubble chart for opens
      this.updateLineChartData();
    } else if (type === 'click') {
      this.stat_type = 'click';
      this.updateBarChartClick(); // Call function to update bar chart for clicks
      this.updateBubbleChartClick(); // Call function to update bubble chart for clicks
      this.updateLineChartData();
      this.updateScatterChartData();
    } else {
      this.stat_type = 'all';
      this.updateBarChartOpen(); // Update both charts for 'all'
      this.updateBarChartClick(); 
      this.updateBubbleChartOpen(); 
      this.updateBubbleChartClick(); 
      this.updateLineChartData();
      this.updateScatterChartData();
    }

    this.reorder = true; // Show the charts again
    this.cdr.detectChanges(); // Trigger change detection again
  }, 0); // Ensure this runs after DOM updates
}


downloadPDF() {
  const pdf = new jsPDF('p', 'mm', 'a4');
  let position = 10; // Starting position for the first chart
  const imgWidth = 190; // Width of the image in the PDF
  const imgHeight = 108; // Height of the image in the PDF

  const chartIds = [
    'barChart',
    'pieChart',
    'radialChart',
    'lineChart',
    'scatterChart',
    'barChartOpen',
    'barChartClick',
    'bubbleChartOpen',
    'bubbleChartClick'
  ];

  const addChartToPDF = async (chartId: string) => {
    const element = document.getElementById(chartId);

    if (!element) {
      console.error(`Element with ID ${chartId} not found.`);
      return;
    }

    try {
      const canvas = await html2canvas(element, {
        scale: 1, // Reduce scale to improve speed and reduce file size
        useCORS: true // Enable if you have external images
      });
      const imgData = canvas.toDataURL('image/png');

      if (position + imgHeight > 297) { // 297mm is the height of an A4 page
        pdf.addPage();
        position = 10; // Reset position for the new page
      }

      pdf.addImage(imgData, 'PNG', 10, position, imgWidth, imgHeight);
      position += imgHeight + 10; // Update position for the next chart
    } catch (error) {
      console.error(`Error capturing chart ${chartId}:`, error);
    }
  };

  const processCharts = async () => {
    for (const chartId of chartIds) {
      await addChartToPDF(chartId);
    }
    pdf.save('charts.pdf');
    console.log('PDF generated successfully.');
  };

  processCharts().catch(error => {
    console.error('Error generating PDF:', error);
  });
}
}
