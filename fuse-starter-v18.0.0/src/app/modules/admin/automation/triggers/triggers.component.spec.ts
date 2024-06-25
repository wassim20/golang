import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TriggersComponent } from './triggers.component';

describe('TriggersComponent', () => {
  let component: TriggersComponent;
  let fixture: ComponentFixture<TriggersComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [TriggersComponent]
    });
    fixture = TestBed.createComponent(TriggersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
