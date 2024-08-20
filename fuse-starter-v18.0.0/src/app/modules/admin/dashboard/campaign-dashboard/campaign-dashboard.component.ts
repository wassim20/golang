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
  imports: [MatFormFieldModule, MatSelectModule,RouterModule,
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
  public scatterChartOptionsOpen: Partial<ScatterChartOptions>;
  public scatterChartOptionsClick: Partial<ScatterChartOptions>;

  campaigns: any;
  campaignId: any = null;
  contacts: any;
  selectedCampaign: any;
  logs: any;
  isLoading: boolean = false;
 

   constructor(private service : DashboardService,private cdr: ChangeDetectorRef,private router: Router,private route:ActivatedRoute) {  
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

    // Initialize pie chart options
    this.pieChartOptions = {
      series: [0, 0],
      chart: {
        width: 500,
        type: "pie"
      },
      labels: ["Total Contacts ;", "Opened Emails", "Clicked Emails"],
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

    this.scatterChartOptionsOpen = {
      series: [{
        name: "Opens",
        data: []
      }],
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
        tickAmount: 12,
        labels: {
          formatter: function (val) {
            return `${Math.floor(val+1)}`;
          }
        }
      },
      tooltip: {
        enabled: true,
        x: {
          formatter: (val) => `Day: ${this.getDayName(val)}`
        },
        z: {
          formatter: (val) => ` ${val}`
        },
      },
      title: {
        text: 'Email Opens by Time of Day',
        align: 'center'
      },
     
    };
    
    this.scatterChartOptionsClick = {
      series: [{
        name: "Clicks",
        data: []
      }],
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
        tickAmount: 24,
        labels: {
          formatter: function (val) {
            return `${Math.floor(val+1)}`;
          }
        }
      },
      tooltip: {
        enabled: true,
        x: {
          formatter: (val) => `Day: ${this.getDayName(val)}`
        },
        y: {
          formatter: (val) => `Hour: ${val}`
        },
        z: {
          formatter: (val) => `${val} times`
        },
      },
      title: {
        text: 'Email Clicks by Time of Day',
        align: 'center'
      },
      
    };
    
    
  }
  
  
  
  
  ngOnInit(): void {
    this.isLoading = true;
    this.route.paramMap.subscribe(params => {
      this.campaignId = params.get('campaignID');
      console.log("campaignID", this.campaignId);
      
    });
    
    this.service.getCampaigns(1, 10).subscribe({
      next: (data) => {
        if (data && data.data && data.data.items) {
          this.campaigns = data.data.items; // Assign mailing list items
          // Select the first campaign by default if not already selected
          
        } else {
          console.error('Invalid response structure:', data);
        }
      },
      error: (error) => {
        console.error('Error fetching campaigns:', error);
      }
    });
    if (this.campaignId) {
      this.service.getCampaign(this.campaignId).subscribe({
        next: (data) => {
          if (data && data.data) {
            this.selectedCampaign = data.data;
            
  
            // Convert fetchTrackingData and fetchContacts to return Promises
            const fetchTrackingDataPromise = this.fetchTrackingData(this.selectedCampaign.id);
            const fetchContactsPromise = this.fetchContacts(this.selectedCampaign.mailingListId);
  
            // Use Promise.all to wait for both promises to resolve
            Promise.all([fetchTrackingDataPromise, fetchContactsPromise])
              .then(() => {
                console.log('Both fetchTrackingData and fetchContacts are done.');
              })
              .catch((error) => {
                console.error('Error in one of the fetch operations:', error);
              });
          } else {
            console.error('Invalid response structure:', data);
          }
        },
        error: (error) => {
          console.error('Error fetching campaign:', error);
        }
      });
  
    }
  }
 
   onCampaignChange(campaign) {
    this.selectedCampaign = campaign;
    this.logs = null; // Reset logs data when a new campaign is selected
    this.contacts = null; // Reset contacts data when a new campaign is selected
    this.isLoading = true; // Set isLoading to true when a new campaign is selected
    
    this.fetchTrackingData(campaign.id); // Fetch tracking data for the selected campaign (assuming 'id' property)
    this.fetchContacts(campaign.mailingListId); // Fetch contacts for the selected campaign (assuming 'mailingListID' property)
    this.router.navigate(['/dashboard', campaign.id]);
  
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
                    this.updateScatterChartOpen(),
                    this.updateScatterChartClick()
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
  updateScatterChartClick() {
  const clicksPerHourDay = {};

  this.logs.forEach(log => {
    if (log.clickedAt && log.clickedAt !== "0001-01-01T01:00:00+01:00") {
      const date = new Date(log.clickedAt);
      const hours = date.getHours() + (date.getMinutes() >= 30 ? 1 : 0);
      const dayOfWeek = date.getDay();

      const key = `${dayOfWeek}-${hours}`;
      if (!clicksPerHourDay[key]) {
        clicksPerHourDay[key] = 0;
      }
      clicksPerHourDay[key]++;
    }
  });

  const seriesData = Object.entries(clicksPerHourDay).map(([key, value]) => {
    const [dayOfWeek, hour] = key.split('-').map(Number);
    return {
      x: dayOfWeek,
      y: hour,
      z: value
    };
  });

  this.scatterChartOptionsClick.series = [{
    name: 'Clicks',
    data: seriesData
  }];

  this.cdr.detectChanges();
}
  updateScatterChartOpen() {
  const opensPerHourDay = {};

  this.logs.forEach(log => {
    if (log.openedAt && log.openedAt !== "0001-01-01T01:00:00+01:00") {
      const date = new Date(log.openedAt);
      const hours = date.getHours() + (date.getMinutes() >= 30 ? 1 : 0);
      const dayOfWeek = date.getDay();

      const key = `${dayOfWeek}-${hours}`;
      if (!opensPerHourDay[key]) {
        opensPerHourDay[key] = 0;
      }
      opensPerHourDay[key]++;
    }
  });

  const seriesData = Object.entries(opensPerHourDay).map(([key, value]) => {
    const [dayOfWeek, hour] = key.split('-').map(Number);
    return {
      x: dayOfWeek,
      y: hour,
      z: value
    };
  });

  this.scatterChartOptionsOpen.series = [{
    name: 'Opens',
    data: seriesData
  }];

  this.cdr.detectChanges();
}
  updateBarChartClick() {
    if (!this.logs) {
      console.error('Logs data is not available.');
      return;
    }

    const clickPerDay = Array(7).fill(0);

    this.logs.forEach(log => {
      if (log.clickedAt && log.clickedAt !== "0001-01-01T01:00:00+01:00") {
          const date = new Date(log.clickedAt);
          const dayOfWeek = date.getDay(); // getDay() returns the day of the week (0 for Sunday, 1 for Monday, etc.)
          clickPerDay[dayOfWeek]++;
      }
      this.barChartOptionsClick.series[0].data = clickPerDay;
      
  });

    
  }
  updateBarChartOpen() {
    if (!this.logs) {
      console.error('Logs data is not available.');
      return;
    }
    // Initialize an array with 7 elements for each day of the week
    const opensPerDay = Array(7).fill(0);

    this.logs.forEach(log => {
        if (log.openedAt && log.openedAt !== "0001-01-01T01:00:00+01:00") {
            const date = new Date(log.openedAt);
            const dayOfWeek = date.getDay(); // getDay() returns the day of the week (0 for Sunday, 1 for Monday, etc.)
            opensPerDay[dayOfWeek]++;
        }
    });

    // Update the bar chart data
    this.barChartOptionsOpen.series[0].data = opensPerDay;
}
  updateScatterChartData() {
    const openedData = this.logs.filter(log => log.openedAt && log.openedAt !== "0001-01-01T01:00:00+01:00" && log.openedAt !== "0001-01-01T00:00:00Z");
    const clickedData = this.logs.filter(log => log.clickedAt && log.clickedAt !== "0001-01-01T01:00:00+01:00");

    // Update the scatter chart series data
    this.scatterChartOptions.series = [
      {
        name: "Opened Emails",
        data: openedData.map(log => ({
          x: new Date(log.openedAt).getTime(),
          y: log.clickCount,
          recipientEmail: log.recipientEmail
        }))
      },
      // {
      //   name: "Clicked Emails",
      //   data: clickedData.map(log => ({
      //     x: new Date(log.clickedAt).getTime(),
      //     y: log.clickCount
      //   }))
      // }
    ];
    this.scatterChartOptions.tooltip = {
      enabled: true,
      custom: ({ seriesIndex, dataPointIndex, w }) => {
        // Access the custom data for the tooltip
        const recipientEmail = w.config.series[seriesIndex].data[dataPointIndex].recipientEmail;
        const clickCount = w.config.series[seriesIndex].data[dataPointIndex].y;
        return `<div class="apexcharts-tooltip-title">Recipient Email: ${recipientEmail}</div>
                <div class="apexcharts-tooltip-title">Click Count: ${clickCount}</div>`;
      }
    };
  }
  updateRadialChartData() {
    const totalLogs = this.logs.length;

    // Calculate the number of opened emails
    const openedLogs = this.logs.filter(log => log.openedAt && log.openedAt !== "0001-01-01T00:00:00Z").length;
    // Calculate the percentage of opened emails
    const openedPercentage = totalLogs > 0 ? (openedLogs / totalLogs) * 100 : 0;
  
    this.radialChartOptions.tooltip = {
      enabled: true,
      y: {
              formatter: () => {
                
               
                return `Opened Emails: ${openedLogs}`;
              },
            },
      
  };
    // Update the radial chart with the calculated percentage
    this.radialChartOptions.series = [openedPercentage];
  }
  updateLineChartData() {
  // Helper function to aggregate data by date
  const aggregateDataByDate = (logs, key) => {
    if (!logs) {
      console.error("Logs data is null or undefined");
      return {};
    }
    return logs.reduce((acc, log) => {
      if (log[key] && log[key] !== "0001-01-01T01:00:00+01:00") {
        const date = new Date(log[key]).toISOString().substring(0, 10);
        acc[date] = (acc[date] || 0) + 1;
      }
      return acc;
    }, {});
  };

  const openedData = aggregateDataByDate(this.logs, 'openedAt');
  const clickedData = aggregateDataByDate(this.logs, 'clickedAt');

  // Combine all unique dates from both opened and clicked data
  const allDates = Array.from(new Set([
    ...Object.keys(openedData),
    ...Object.keys(clickedData)
  ])).sort();

  // Create the series data arrays
  const openedSeriesData = allDates.map(date => openedData[date] || 0);
  const clickedSeriesData = allDates.map(date => clickedData[date] || 0);

  // Update the chart options
  this.lineChartOptions.xaxis.categories = allDates;
  this.lineChartOptions.series = [
    {
      name: "Opened Emails",
      data: openedSeriesData
    },
    {
      name: "Clicked Emails",
      data: clickedSeriesData
    }
  ];
}
  updateChartData(): void {
  let totalOpened = 0;
  let totalClicked = 0;
  let totalErrors = 0;

  if (this.logs && this.logs.length > 0) {

  this.logs.forEach(log => {
    if (log.openedAt && new Date(log.openedAt).getFullYear() > 1) {
      totalOpened++;
    }
    if (log.clickCount > 0) {
      totalClicked++;
    }
    if (log.error) {
      totalErrors++;
    }
  });} else{
    console.log('No logs to display.');

  }

  // Adjusting the structure to match a single dataset for all categories
  this.barChartOptions.series = [
    {
      name: "",
      data: [totalOpened,totalClicked,totalErrors]
    },
    
  ];
  
  
}
  updatePieChartData(): void {
  if (!this.contacts) {
    console.error("Contacts data is null or undefined");
    return;
  }

  if (!this.logs) {
    console.error("Logs data is null or undefined");
    return;
  }

  const totalContacts = this.contacts.length;
  const openedEmails = this.logs.filter(log => log.openedAt && new Date(log.openedAt).getFullYear() > 1).length;
  const clickedEmails = this.logs.filter(log => log.clickCount > 0).length;

  this.pieChartOptions.labels = ["Opened Emails", "Clicked Emails", "Total Contacts: " + totalContacts];
  // Update the entire dataset to avoid potential issues
  this.pieChartOptions.series = [openedEmails, clickedEmails];
}
getDayName(dayIndex) {
  const days = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
  return days[dayIndex];
}
}
