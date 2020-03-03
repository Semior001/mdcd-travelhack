import {Component, OnDestroy, OnInit} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {HttpClient} from '@angular/common/http';
import {Subscription} from 'rxjs';
import {environment} from '../../environments/environment';
import {ActivatedRoute, Router} from '@angular/router';
import {MatSnackBar} from '@angular/material';


@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.css']
})
export class UserComponent implements OnInit, OnDestroy {
  isLoading = true;
  userSub: Subscription;
  form: FormGroup;
  user = null;

  constructor(
    private http: HttpClient,
    private route: ActivatedRoute,
    private router: Router,
    private matSnackBar: MatSnackBar
  ) {
  }

  ngOnInit() {
    const userId = Number(this.route.snapshot.paramMap.get('id'));
    if (isNaN(userId)) {
      console.log('Wrong user id');
      this.router.navigate(['main']);
    } else {
      this.userSub = this.http.get(
        environment.apiUrl + '/users/' + userId
        ,
        {withCredentials: true}).subscribe(
        (response: any) => {
          this.user = response;
          console.log(response);
          this.isLoading = false;
          this.form = new FormGroup({
            // TODO state
            'email': new FormControl(this.user.Email, [Validators.required, Validators.email]),
            'password': new FormControl(null, [Validators.required, Validators.minLength(6)])
          });
        },
        (error) => {
          this.matSnackBar.open('Ошибка при получении данных пользователя');
        }
      );
    }
  }

  deleteUser() {
    this.http.delete(
      environment.apiUrl + '/users/' + this.user.ID
      ,
      {withCredentials: true}).subscribe(
      (response: any) => {
        console.log(response);
        this.matSnackBar.open('Пользователь удалён');
      },
      (error) => {
        this.matSnackBar.open('Ошибка при удалении пользователя');
      }
    );
  }

  onSubmit() {
    if (this.form.valid) {
      const data = {email: this.form.value['email'].trim(), password: this.form.value['password'].trim()};
      this.isLoading = true;
      this.http.put(
        environment.apiUrl + `/users/${this.user.ID}`, data
        ,
        {withCredentials: true}).subscribe(
        (response: any) => {
          this.matSnackBar.open('Пользователь обновлён');
          this.isLoading = false;
          this.user.Email = data.email;
          this.form.patchValue({
            'email': data.email,
            'password': null
          });
        },
        (error) => {
          this.matSnackBar.open('Ошибка при сохранении данных пользователя');
        }
      );
    }
  }

  ngOnDestroy(): void {
    if (this.userSub) {
      this.userSub.unsubscribe();
    }
  }
}
