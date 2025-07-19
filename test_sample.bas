10 REM Sample BASIC program to test the interpreter
20 PRINT "=== BASIC Interpreter Test ==="
30 PRINT
40 LET A = 10
50 LET B = 5
60 PRINT "A ="; A
70 PRINT "B ="; B
80 PRINT "A + B ="; A + B
90 PRINT "A * B ="; A * B
100 PRINT
110 PRINT "Counting with FOR loop:"
120 FOR I = 1 TO 5
130 PRINT "  "; I
140 NEXT I
150 PRINT
160 PRINT "Testing IF statements:"
170 IF A > B THEN PRINT "A is greater than B"
180 IF B > A THEN PRINT "B is greater than A"
190 PRINT
200 PRINT "Program completed successfully!"
210 END