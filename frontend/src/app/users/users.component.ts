import {Component, OnDestroy, OnInit} from '@angular/core';
import {environment} from '../../environments/environment';
import {Subscription} from 'rxjs';
import {HttpClient} from '@angular/common/http';
import {ActivatedRoute} from '@angular/router';
import {MatSnackBar} from '@angular/material';
import {FormControl, FormGroup, Validators} from '@angular/forms';

@Component({
  selector: 'app-users',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.css']
})
export class UsersComponent implements OnInit, OnDestroy {
  users = [];
  isLoading = true;
  usersSub: Subscription;
  form: FormGroup;

  constructor(
    private http: HttpClient,
    private  matSnackBar: MatSnackBar
  ) {
  }

  onSubmit() {
    const data = {
      'email': this.form.value['email'].trim(),
      'password': this.form.value['password'].trim(),
      'priveleges': {
        'admin': this.form.value['role'] === '1'
      }
    };
    this.http.post(
      environment.apiUrl + '/users',
      data,
      {withCredentials: true}).subscribe(
      (response: any[]) => {
        this.matSnackBar.open('Пользователь создан');
        const now = Date.now();
        this.users.push({
          ID: response['ID'],
          Email: data.email,
          UpdatedAt: now,
          CreatedAt: now,
          Priveleges: data.priveleges
        });
      },
      (error) => {
        this.matSnackBar.open('Ошибка при создании пользователя');
      });
  }

  ngOnInit() {
    this.form = new FormGroup({
      'email': new FormControl(null, [Validators.required, Validators.email]),
      'password': new FormControl(null, [Validators.required, Validators.minLength(5)]),
      'role': new FormControl(null, [Validators.required])
    });

    this.usersSub = this.http.get(
      environment.apiUrl + '/users',
      {withCredentials: true}).subscribe(
      (response: any[]) => {
        this.users = response;
        this.isLoading = false;
      },
      (error) => {
        this.matSnackBar.open('Ошибка при получении пользователей');
      }
    );
  }

  ngOnDestroy() {
    this.usersSub.unsubscribe();
  }
}
