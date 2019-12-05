func line ( int x , in y ){
  // last number in loop is the step
  iter (i := 0; x , 1){ // declares it , scope is the loop
    circle (2 , 3, y , 5);
    }
}

// macro entry
func main (){
  iter (i := 0; 3 , 1){
    rect (i , i , 3 , 0 xff );
  }

  iter (j := 0; 8 , 2){ // loops 0 2 4 6 8
    rect (j , j , 8 , 0 xff );
  }

  circle (4 , 5, 2, 0 x11000011 );
}
