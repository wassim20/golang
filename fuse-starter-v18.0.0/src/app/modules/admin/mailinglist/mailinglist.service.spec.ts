import { TestBed } from '@angular/core/testing';

import { MailinglistService } from './mailinglist.service';

describe('MailinglistService', () => {
  let service: MailinglistService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(MailinglistService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
