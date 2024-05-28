import { Component, ViewEncapsulation, ElementRef, OnInit, ViewChild, AfterViewInit } from '@angular/core';
import { EmailEditorComponent,EmailEditorService,UnlayerOptions,EmailEditorModule  } from '@trippete/angular-email-editor';
import {MatButtonModule} from '@angular/material/button';


@Component({
    selector     : 'example',
    standalone   : true,
    imports     : [EmailEditorModule,MatButtonModule],
    templateUrl  : './example.component.html',
    encapsulation: ViewEncapsulation.None,
})
export class ExampleComponent
{
    
  @ViewChild(EmailEditorComponent)
  private emailEditor!: any;

  theme : string = 'dark';

  unlayerOptions!: UnlayerOptions;

  exportHtml() {
    this.emailEditor.exportHtml((data :any ) => {
      const { design, html } = data;
      console.log('exportHtml', html);
    }, {
      cleanup: true,
    });
  }

  saveDesign() {
    this.emailEditor.saveDesign((design :any) => {
      console.log('saveDesign', design);
    });
  }

  loadDesign() {
    const design = {}; // Your saved design JSON
    this.emailEditor.loadDesign(design);
  }

  onEditorLoad() {
    // Handle any actions after the editor loads (optional)
  }
  


   
} 
