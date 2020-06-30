func line ( int x , int y ){
  iter (i := 0; x , 1){
    circle (2 , 3, y , 5);
  }
}

func gg (){
  iter (i := 0; 3 , 1){
    rect (i , i , 3 , 0xf);
  }

  iter (j := 0; 8 , 2){ // loops 0 2 4 6 8
    rect (j , j , 8 , 0xf);
  }

  circle (4 , 5,  2 , 0xf);
}
