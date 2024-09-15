import { AfterViewChecked, AfterViewInit, ChangeDetectionStrategy, ChangeDetectorRef, Component, OnInit, ViewChild, ViewEncapsulation } from '@angular/core';
import { DashboardService } from './dashboard.service';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { CommonModule } from '@angular/common'; // Import CommonModule for *ngFor
import { MatOptionModule } from '@angular/material/core'; // Import MatOptionModule for mat-option
import { HttpClientModule } from '@angular/common/http';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import {MatDatepickerModule} from '@angular/material/datepicker';
import jsPDF from 'jspdf';

import ApexCharts from 'apexcharts';


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
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIcon, MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatTooltipModule } from '@angular/material/tooltip';
import html2canvas from 'html2canvas';
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
  selector: 'dashboard',
  standalone: true,
  imports: [MatFormFieldModule, MatSelectModule,RouterModule,MatDatepickerModule, FormsModule, ReactiveFormsModule,MatIconModule,
    MatButtonModule,MatTooltipModule,NgApexchartsModule,
    CommonModule, MatOptionModule, HttpClientModule,MatProgressSpinnerModule,NgApexchartsModule],
  templateUrl: './dashboard.component.html',
  providers: [],
  changeDetection: ChangeDetectionStrategy.OnPush,
  encapsulation: ViewEncapsulation.None,
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  public barChartOptions: Partial<BarChartOptions>;
  public barChartOptionsOpen: Partial<ColumnChartOptions>;
  public barChartOptionsClick: Partial<ColumnChartOptions>;
  public pieChartOptions: Partial<PieChartOptions>;
  public lineChartOptions: Partial<LineChartOptions>;
  public radialChartOptions: Partial<RadialChartOptions>;
  public scatterChartOptions: Partial<ScatterChartOptions>;
  public bubbleChartOptionsOpen: Partial<BubbleChartOptions>;
  public bubbleChartOptioncClick: Partial<BubbleChartOptions>;
  // Accessing chart components using @ViewChild
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
  stat_type: any = "all";
  reorder : boolean = true;
  campaignId: any = null;
  contacts: any;
  OpenRatio: number;
  ClickRatio: number;
  selectedCampaign: any;
  logs: any;
  isLoading: boolean = false;
  selectedStat: string = 'all'; // Default to 'all'

 

   constructor(private route :ActivatedRoute,private service : DashboardService,private cdr: ChangeDetectorRef,private router: Router) {  
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

  // Fetch campaigns list
  this.service.getCampaigns(1, 10).subscribe({
    next: async (data) => {
      if (data && data.data && data.data.items) {
        this.campaigns = data.data.items;

        // Check if "all" is selected
        if (this.route.snapshot.paramMap.get('campaignID') === null) {
          this.selectedCampaign = 'all';
        } else {
          this.selectedCampaign = this.campaigns.find(campaign => campaign.id === this.route.snapshot.paramMap.get('campaignID'));
        }

        try {
          await this.fetchAllContacts();
          await this.fetchAllTrackingData();
        } catch (error) {
          console.error('Error fetching data:', error);
        } finally {
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


 
  onCampaignChange(campaign) {
    if (campaign === 'all') {
      this.selectedCampaign = 'all';
      this.router.navigate(['/dashboard']);
    } else {
      this.selectedCampaign = campaign;
      this.logs = null;
      this.contacts = null;
  
      this.router.navigate(['/dashboard', campaign.id]).then(() => {
        this.selectedCampaign = this.campaigns.find(c => c.id === campaign.id);
      });
    }
  }
  fetchAllContacts() {
    this.service.getAllContacts().subscribe(
        (data) => {
            this.contacts = data?.data?.items?.length > 0 ? data.data.items : null;
            
            if(this.contacts && this.contacts.length > 0) {
            //this.updatePieChartData();
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
  fetchAllTrackingData() {
    this.service.getAllTrackingLogs().subscribe(
        (data) => {
            this.logs = data?.data?.items?.length > 0 ? data.data.items : null;
            
            // if logs are not empty 
            if (this.logs && this.logs.length > 0) {
                const promises = [
                    this.updateChartData(),
                    this.updatePieChartData(),
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
                }).catch((error) => {
                    console.error('Error updating charts:', error);
                });
            } else {
                console.log('No logs to display.');
            }
          },
          (error) => {
            console.error('Error fetching logs:', error);
            this.logs = null; // Optionally, also set logs to null in error scenario
          }
        );
      }
      updateChartData(): void {
        let totalOpened = 0;
        let totalClicked = 0;
        let totalErrors = 0;
        
        this.service.updateChartData().subscribe((data) => {
          
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
        this.cdr.detectChanges();
        
        
      }
      updateBubbleChartClick(): void {
        this.service.bubbleChartDataClicks().subscribe((response) => {
          
      
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
          this.cdr.detectChanges();
      
          // Trigger change detection to update the chart
        },
        (error) => {
          console.error('Error fetching bubble chart click data:', error);
          // Handle the error appropriately, e.g., display an error message
        });
      }
      
      updateBubbleChartOpen(): void {
        this.service.bubbleChartDataOpens().subscribe((response) => {
          
      
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
          this.cdr.detectChanges();
        },
        (error) => {
          console.error('Error fetching bubble chart data:', error);
          // Handle the error appropriately, e.g., display an error message
        });
        this.cdr.detectChanges();
      }
      
  updateBarChartClick() {
   this.service.barChartDataClicks().subscribe((data) => {
      this.barChartOptionsClick.series[0].data = data.data;
      this.cdr.detectChanges();      
    });
    this.cdr.detectChanges();      

    
  }
  updateBarChartOpen() {
    this.service.barChartDataOpens().subscribe((data) => {

    

    // Update the bar chart data
    this.barChartOptionsOpen.series[0].data = data.data;
    this.cdr.detectChanges();
  });
  this.cdr.detectChanges();
}
updateScatterChartData(): void {
  this.service.updateScatterChartData().subscribe((data) => {

    // Update the scatter chart series data using the data from the backend
    this.scatterChartOptions.series = [
      
      {
        name: "Opened Emails",
        data: data.data.openedData.map(item => ({
          
          x: item.x,  // Timestamp already provided by the backend
          y: item.y,  // Click count provided by the backend
          recipientEmail: item.recipientEmail
        }))
      },
      {
        name: "Clicked Emails",
        data: data.data.clickedData.map(item => ({
          x: item.x,  // Timestamp already provided by the backend
          y: item.y,  // Click count provided by the backend
          recipientEmail: item.recipientEmail
        }))
      }
    ];
    this.cdr.detectChanges();
    
    // Update tooltip to use the custom data
    this.scatterChartOptions.tooltip = {
      enabled: true,
      custom: ({ seriesIndex, dataPointIndex, w }) => {
        const recipientEmail = w.config.series[seriesIndex].data[dataPointIndex].recipientEmail;
        const clickCount = w.config.series[seriesIndex].data[dataPointIndex].y;
        return `<div class="apexcharts-tooltip-title">Recipient Email: ${recipientEmail}</div>
        <div class="apexcharts-tooltip-title">Click Count: ${clickCount}</div>`;
      }
    };
    
    // Trigger change detection to update the chart this.cdr.detectChanges();
    this.cdr.detectChanges();
  },
  (error) => {
    console.error('Error fetching scatter chart data:', error);
    // Optionally handle the error, such as displaying a notification to the user
  });
}

  updateRadialChartData() {
    const totalLogs = this.logs.length;

    this.service.updateRadialChartData().subscribe((data) => {


      
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
      this.cdr.detectChanges();
    });
    this.cdr.detectChanges();
    }
    updateLineChartData(): void {
      this.service.updateLineChartData().subscribe((data) => {
        if (this.lineChartOptions?.chart) {
          const allDates = data.data.allDates.map(date => `${date}`);
          const totals = data.data.totals;
    
          this.lineChartOptions = {
            ...this.lineChartOptions,
            xaxis: {
              categories: allDates,
            },
            series: [
              { name: "Opened Emails", data: data.data.openedSeriesData },
              { name: "Clicked Emails", data: data.data.clickedSeriesData }
            ],
            title: {
              text: `Email Campaign Tracking (Total Sent: ${totals})`,
              align: "left"
            }
          };
          
          this.cdr.detectChanges();  // Only detect changes once at the end
        } else {
          console.error("Line chart is not initialized");
        }
      });
    }
    
    
    updatePieChartData(): void {
      
    
      const totalContacts = this.contacts.length;
    
      this.service.updatePieChartData().subscribe(
        (data) => {
          
          const openedEmails = data.data.opened_count || 0;
          const clickedEmails = data.data.clicked_count || 0;
    
          this.pieChartOptions = {
            ...this.pieChartOptions,
            labels: ["Opened Emails", "Clicked Emails", "Total Contacts: " + totalContacts],
            series: [openedEmails, clickedEmails],
          };
          
          this.cdr.detectChanges();  // Only detect changes once at the end
        },
        (error) => {
          console.error('Error fetching data from backend:', error);
        }
      );
    }
    

// Separate method to update chart options
private updateChartOptions(totalContacts: number, openedEmails: number, clickedEmails: number): void {
  this.pieChartOptions = {
    ...this.pieChartOptions,
    labels: ["Opened Emails", "Clicked Emails", `Total Contacts: ${totalContacts}`],
    series: [openedEmails, clickedEmails]
  };

  // Trigger change detection to update the view
  this.cdr.detectChanges();
}
getDayName(dayIndex) {
  const days = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
  return days[dayIndex];
}

//stat change 
statchange(type: string) {
  this.selectedStat = type;
  this.reorder = false; // Temporarily hide the chart
  this.cdr.detectChanges(); // Trigger change detection

  setTimeout(() => {
    if (type === 'open') {
      this.stat_type = 'open';
    } else if (type === 'click') {
      this.stat_type = 'click';
    } else {
      this.stat_type = 'all';
    }
    this.reorder = true; // Show the chart again
    this.cdr.detectChanges(); // Trigger change detection
  }, 0); // Delay to ensure the DOM updates
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




