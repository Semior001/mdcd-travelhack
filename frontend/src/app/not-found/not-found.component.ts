import {Component, OnInit, ViewEncapsulation} from '@angular/core';

@Component({
  selector: 'app-not-found',
  template: `<div class="not-found-container">
                <h1>HTTP 404 - NOT FOUND</h1>
             </div>`,
  styleUrls: ['./not-found.component.css'],
  encapsulation: ViewEncapsulation.None
})
export class NotFoundComponent {}
