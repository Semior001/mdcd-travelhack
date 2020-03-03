import {Component, OnInit} from '@angular/core';
import {environment} from '../../environments/environment';
import {HttpClient, HttpEventType} from '@angular/common/http';
import {MatSnackBar} from '@angular/material';

@Component({
  selector: 'app-backgrounds',
  templateUrl: './backgrounds.component.html',
  styleUrls: ['./backgrounds.component.css']
})
export class BackgroundsComponent implements OnInit {
  selectedFile: File;

  constructor(
    private http: HttpClient,
    private matSnackBar: MatSnackBar
  ) {
  }

  ngOnInit() {
  }

  onFileSelected(event) {
    this.selectedFile = <File>event.target.files[0];
    this.onUpload();
  }

  onUpload() {
    const body = new FormData();
    body.append('image', this.selectedFile);
    this.http.post(
      environment.apiUrl + '/add_image?imgType=background',
      body, {
        reportProgress: true,
        observe: 'events',
        withCredentials: true
      }
    ).subscribe((event) => {
        if (event.type === HttpEventType.UploadProgress) {
          console.log('Upload progress: ' + Math.round(event.loaded / event.total * 100) + '%');
        } else if (event.type === HttpEventType.Response) {
          console.log(event);
          this.matSnackBar.open('Изображение загружено');
        }
        console.log(event);
      },
      (error) => {
        this.matSnackBar.open('Ошибка при загрузке фотографии');
      }
    );
    console.log('Upload background');
  }
}
