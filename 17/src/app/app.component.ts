import { Component,ViewChild } from '@angular/core';
import { EmailEditorComponent,EmailEditorService,UnlayerOptions  } from '@trippete/angular-email-editor';
declare var unlayer: any;

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent {
  title = '17';
  customTool = {
    name: 'custom_tool', // Unique identifier
    label: 'Custom Tool', // Display label
    icon: 'fa fa-star', // Icon
    supportedDisplayModes: ['email', 'web'], // Supported display modes
    usageLimit: 5, // Usage limit
    renderer: {
      type: 'text', // Type of content
      value: 'Hello World' // Content value
    },
    properties: {
      // Define your properties here
    }
  };
  tableTool = {
    name: 'table_tool', // Unique identifier
    label: 'Table Tool', // Display label
    icon: 'fa fa-table', // Icon
    supportedDisplayModes: ['email', 'web'], // Supported display modes
    usageLimit: 10, // Usage limit
    renderer: {
      type: 'html', // Type of content
      value: `
        <table>
          <tr>
            <th>Header 1</th>
            <th>Header 2</th>
          </tr>
          <tr>
            <td>Cell 1</td>
            <td>Cell 2</td>
          </tr>
        </table>
      ` // Content value
    },
    properties: {
      // Define your properties here
    }
  };

 

  @ViewChild(EmailEditorComponent)
  private emailEditor!: EmailEditorComponent;

  unlayerOptions!: UnlayerOptions;

  ngAfterViewInit() {
    this.emailEditor.loaded.subscribe(() => {
      console.log('Email editor loaded!');
      this.registerCustomTools(); // Call your function to register tools
    });
  }
  registerCustomTools() {
    if (this.emailEditor && this.emailEditor.editor) {
    this.emailEditor.editor.registerTool({
      name: 'table',
      label: 'Table',
      icon: 'fa-table',
      supportedDisplayModes: ['email'],
      renderer: {
        customJS: unlayer.createViewer({
          render(values: { rows: any[]; }) {
            const rows = values.rows.map((row: { cells: any[]; }) => (
              `<tr>
                ${row.cells.map((cell: { content: any; }) => `<td>${cell.content}</td>`).join('')}
              </tr>`
            )).join('');
            return `<table><tbody>${rows}</tbody></table>`;
          }
        }),
        exporters: {
          email: function (values: { rows: any[]; }) {
            return values.rows.map((row: { cells: any[]; }) => (
              `<tr>
                ${row.cells.map((cell: { content: any; }) => `<td>${cell.content}</td>`).join('')}
              </tr>`
            )).join('');
          }
        },
        head: {
          css: function (values: any) {
            return `
              table {
                border-collapse: collapse;
              }
              th, td {
                border: 1px solid #ddd;
                padding: 8px;
              }
            `;
          },
          js: function (values: any) {
            return ''; // No additional JavaScript needed
          }
        }
      },
      options: {
        rows: {
          label: 'Rows',
          widget: 'array',
          value: [
            {
              cells: [
                { content: '' },
                { content: '' },
                { content: '' }
              ]
            }
          ],
          schema: [
            {
              name: 'cells',
              label: 'Cells',
              widget: 'array',
              value: [],
              schema: [
                { name: 'content', label: 'Content', widget: 'text' }
              ]
            }
          ]
        }
      },
      values: {
        rows: [
          { cells: [{ content: '' }, { content: '' }, { content: '' }] }
        ]
      },
      validator(data: any) {
        return [];
      }
    });
    this.emailEditor.editor.registerTool({
      
        name: 'social_media_icons',
        label: 'Social Media Icons',
        icon: 'fa-share-alt',
        supportedDisplayModes: ['email'],
        renderer: {
          customJS: unlayer.createViewer({
            render(values: { socialLinks: any[]; }) {
              const iconUrls: { [key: string]: string } = {
                facebook: 'path/to/facebook_icon.png',
                // Add URLs for other social media icons
              };
              return values.socialLinks.map((link: { url: any; platform: string | number; }) => (
                `<a href="${link.url}" target="_blank"><img src="${iconUrls[link.platform]}" /></a>`
              )).join('');
            }
          }),
          exporters: {
            email: function (values: { socialLinks: any[]; }) {
              return values.socialLinks.map((link: { url: any; platform: any; }) => (
                `<a href="${link.url}" target="_blank"><img src="${link.platform}_icon.png" /></a>`
              )).join('');
            }
          },
          head: {
            css: function (values: any) {
              return '.social-icon { width: 20px; height: 20px; }';
            },
            js: function (values: any) {
              return ''; // No additional JavaScript needed
            }
          }
        },
        options: {
          socialLinks: {
            label: 'Social Links',
            widget: 'array',
            value: [],
            schema: [
              {
                name: 'platform',
                label: 'Platform',
                widget: 'select',
                options: [
                  { value: 'facebook', label: 'Facebook' },
                  // Add options for other platforms
                ]
              },
              {
                name: 'url',
                label: 'URL',
                widget: 'text',
              }
            ]
          }
        },
        values: {
          socialLinks: []
        },
        validator(data: any) {
          return [];
        }
      });
    
    }else {
      console.log('Email editor not fully initialized yet');
    }
  }
  
 
 
  
  
  

  // called when the editor is created
  editorLoaded() {
    console.log('editorLoaded');
   
   this.emailEditor.editor.loadDesign({
   });
  }

  // called when the editor has finished loading
  editorReady() {
    console.log('editorReady');
  }

  exportHtml() {
    this.emailEditor.editor.exportHtml((data: any) =>
      console.log('exportHtml', data.html)
    );
  }

}
