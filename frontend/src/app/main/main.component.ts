import {Component, OnInit} from '@angular/core';
import {HttpClient, HttpEventType} from '@angular/common/http';
import {environment} from '../../environments/environment';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {MatSnackBar} from '@angular/material';
import {DomSanitizer} from '@angular/platform-browser';

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.css']
})
export class MainComponent implements OnInit {
  isLoading = true;
  imageId: number = null;
  filterName: string = null;
  backgroundId: number = null;
  imageUrl: string = null;
  form: FormGroup;
  selectedFile: File = null;
  imgToUpload: File = null;
  background_ids = [];
  apiUrl = environment.apiUrl;
  filters: string[] = [
    '_1977',
    'aden',
    'brannan',
    'brooklyn',
    'clarendon',
    'css-contrast',
    'css-grayscale',
    'css-hue_rotate',
    'css-saturate',
    'css-sepia',
    'earlybird',
    'gingham',
    'hudson',
    'inkwell',
    'kelvin',
    'lark',
    'lofi',
    'maven',
    'mayfair',
    'moon',
    'nashville',
    'perpetua',
    'reyes',
    'rise',
    'slumber',
    'stinson',
    'toaster',
    'valencia',
    'walden',
    'willow',
    'xpro2',
  ];
  barcode: number = null;

  constructor(
    private http: HttpClient,
    private matSnackBar: MatSnackBar,
    public _DomSanitizationService: DomSanitizer
  ) {
  }

  onBarcodeSubmit() {
  //   this.http.get(
  //     environment.apiUrl + `/get_image_from_barcode?barcode=${this.form.value['barcode'].trim()}`,
  //     {withCredentials: true}
  //   ).subscribe(
  //     (response: Blob) => {
  //       console.log(response);
  //       TODO
        // this.imageUrl = window.URL.createObjectURL(response);
      // },
      // (error) => {
      //   this.matSnackBar.open('Ошибк при получении изображения');
      // }
    // );
  }

  onImgUpload() {
    const body = new FormData();
    body.append('image', this.imgToUpload);
    this.http.post(
      environment.apiUrl + `/add_image?imgType=photo&barcode=${this.form.value['barcode']}`,
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
          this.barcode = this.form.value['barcode'];
          this.imageUrl = window.URL.createObjectURL(this.imgToUpload);
          this.imageId = event.body['ID'];
        }
        console.log(event);
      },
      () => {
        this.matSnackBar.open('Ошибка при загрузке фотографии');
      }
    );
    console.log('Upload background');
  }

  OnImgToUploadSelected(event) {
    this.imgToUpload = <File>event.target.files[0];
    this.onImgUpload();
  }

  ngOnInit() {
    this.form = new FormGroup({
      'barcode': new FormControl(null, [Validators.required])
    });
    this.http.get(environment.apiUrl + '/backgrounds',
      {withCredentials: true}).subscribe(
      (response: any[]) => {
        this.background_ids = response;
        this.isLoading = false;
      },
      (error) => {
        this.matSnackBar.open('Ошибка при получении фоновых изображений');
      }
    );
  }

  onBackgroundClick(backgroundId) {
    if (this.barcode !== null) {
      const body = new FormData();
      this.selectedFile = null;
      this.backgroundId = backgroundId;
      body.append('barcode', this.barcode.toString());
      body.append('background_id', this.backgroundId.toString());
      if (this.filterName === null) {
        body.append('filter_name', 'base');
      } else {
        body.append('filter_name', this.filterName);
      }
      this.http.post(
        environment.apiUrl + '/apply_effects',
        body,
        {withCredentials: true}).subscribe(
        (response: File) => {
          this.imageUrl = window.URL.createObjectURL(response);
        },
        (error) => {
          this.matSnackBar.open('Ошибка при обработке фотографии');
        }
      );
      console.log(backgroundId);
    }
  }

  resetMain() {
    this.filterName = null;
    this.backgroundId = null;
    this.selectedFile = null;
    this.barcode = null;
    this.imageUrl = null;
    this.form.patchValue({
      'barcode': null
    });
  }

  onFilterClick(filterName: string) {
    if (this.barcode !== null) {
      const body = new FormData();
      body.append('barcode', this.barcode.toString());
      body.append('filter_name', filterName);
      if (this.backgroundId === null) {
        if (this.selectedFile === null) {
          body.append('background_id', '-1');
        } else {
          body.append('background', this.selectedFile);
        }
      } else {
        body.append('background_id', this.backgroundId.toString());
      }
      this.http.post(
        environment.apiUrl + '/apply_effects',
        body,
        {withCredentials: true}).subscribe(
        (response: Blob | File) => {
          const URL =  window.URL.createObjectURL || webkitURL.createObjectURL;
          this.imageUrl = URL(response);
        },
        (error) => {
          this.matSnackBar.open('Ошибка при обработке фотографии');
        }
      );
      console.log(filterName);
    }
  }

  onFileSelected(event) {
    this.backgroundId = null;
    this.selectedFile = <File>event.target.files[0];
    this.onUpload();
  }

  onUpload() {
    if (this.barcode !== null) {
      const body = new FormData();
      body.append('barcode', this.barcode.toString());
      body.append('backgound', this.selectedFile);
      if (this.filterName === null) {
        body.append('filter_name', 'base');
      } else {
        body.append('filter_name', this.filterName);
      }
      this.http.post(
        environment.apiUrl + '/apply_effects',
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
            // this.imageUrl = window.URL.createObjectURL(response);
          }
          console.log(event);
        },
        (error) => {
          this.matSnackBar.open('Ошибка при обработке фотографии');
        }
      );
      console.log('Upload background');
    }
  }

  saveChanges() {
    if (this.barcode !== null &&
      (this.filterName !== null || this.backgroundId !== null || this.selectedFile !== null)
    ) {
      console.log('Saving changes');
    }
  }
}
