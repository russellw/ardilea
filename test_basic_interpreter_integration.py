"""Integration tests for BASIC interpreter with line number support."""

import pytest
from io import StringIO
import sys
from basic_interpreter import BasicInterpreter


class TestBasicInterpreterIntegration:
    """Integration tests for the BASIC interpreter."""
    
    def setup_method(self):
        """Set up test fixtures."""
        self.interpreter = BasicInterpreter()
        self.captured_output = StringIO()
        
    def capture_output(self):
        """Capture stdout for testing output."""
        sys.stdout = self.captured_output
        
    def restore_output(self):
        """Restore stdout."""
        sys.stdout = sys.__stdout__
        
    def get_output(self):
        """Get captured output."""
        return self.captured_output.getvalue().strip()
    
    def test_simple_print_statement(self):
        """Test basic PRINT statement execution."""
        program = "10 PRINT \"Hello, World!\""
        
        self.capture_output()
        self.interpreter.run(program)
        self.restore_output()
        
        assert self.get_output() == "Hello, World!"
    
    def test_multiple_line_program(self):
        """Test program with multiple line numbers."""
        program = """10 PRINT "First line"
20 PRINT "Second line"
30 PRINT "Third line" """
        
        self.capture_output()
        self.interpreter.run(program)
        self.restore_output()
        
        expected = "First line\nSecond line\nThird line"
        assert self.get_output() == expected
    
    def test_line_number_ordering(self):
        """Test that lines execute in line number order, not input order."""
        program = """30 PRINT "Third"
10 PRINT "First"
20 PRINT "Second" """
        
        self.capture_output()
        self.interpreter.run(program)
        self.restore_output()
        
        expected = "First\nSecond\nThird"
        assert self.get_output() == expected
    
    def test_variable_assignment_and_usage(self):
        """Test variable assignment and usage."""
        program = """10 LET A = 42
20 PRINT A"""
        
        self.capture_output()
        self.interpreter.run(program)
        self.restore_output()
        
        assert self.get_output() == "42"
    
    def test_arithmetic_operations(self):
        """Test basic arithmetic operations."""
        program = """10 LET A = 10
20 LET B = 5
30 PRINT A + B
40 PRINT A - B
50 PRINT A * B
60 PRINT A / B"""
        
        self.capture_output()
        self.interpreter.run(program)
        self.restore_output()
        
        expected = "15\n5\n50\n2"
        assert self.get_output() == expected
    
    def test_goto_statement(self):
        """Test GOTO statement for program flow control."""
        program = """10 PRINT "First"
20 GOTO 40
30 PRINT "This should not print"
40 PRINT "Last" """
        
        self.capture_output()
        self.interpreter.run(program)
        self.restore_output()
        
        expected = "First\nLast"
        assert self.get_output() == expected
    
    def test_if_statement(self):
        """Test IF statement conditional execution."""
        program = """10 LET A = 10
20 IF A > 5 THEN PRINT "A is greater than 5"
30 IF A < 5 THEN PRINT "A is less than 5"
40 PRINT "Done" """
        
        self.capture_output()
        self.interpreter.run(program)
        self.restore_output()
        
        expected = "A is greater than 5\nDone"
        assert self.get_output() == expected
    
    def test_for_loop(self):
        """Test FOR loop execution."""
        program = """10 FOR I = 1 TO 3
20 PRINT I
30 NEXT I"""
        
        self.capture_output()
        self.interpreter.run(program)
        self.restore_output()
        
        expected = "1\n2\n3"
        assert self.get_output() == expected
    
    def test_nested_for_loops(self):
        """Test nested FOR loops."""
        program = """10 FOR I = 1 TO 2
20 FOR J = 1 TO 2
30 PRINT I; J
40 NEXT J
50 NEXT I"""
        
        self.capture_output()
        self.interpreter.run(program)
        self.restore_output()
        
        expected = "1 1\n1 2\n2 1\n2 2"
        assert self.get_output() == expected
    
    def test_input_statement(self):
        """Test INPUT statement for user input."""
        program = """10 INPUT "Enter a number: "; A
20 PRINT "You entered: "; A"""
        
        # Mock input
        original_input = input
        inputs = iter(["42"])
        input_func = lambda prompt="": next(inputs)
        
        try:
            # Replace input function
            import builtins
            builtins.input = input_func
            
            self.capture_output()
            self.interpreter.run(program)
            self.restore_output()
            
            output = self.get_output()
            assert "Enter a number:" in output
            assert "You entered: 42" in output
        finally:
            # Restore original input
            builtins.input = original_input
    
    def test_string_operations(self):
        """Test string variable operations."""
        program = """10 LET A$ = "Hello"
20 LET B$ = "World"
30 PRINT A$; " "; B$; "!" """
        
        self.capture_output()
        self.interpreter.run(program)
        self.restore_output()
        
        assert self.get_output() == "Hello World!"
    
    def test_line_number_gaps(self):
        """Test that programs work with non-consecutive line numbers."""
        program = """100 PRINT "Line 100"
500 PRINT "Line 500"
1000 PRINT "Line 1000" """
        
        self.capture_output()
        self.interpreter.run(program)
        self.restore_output()
        
        expected = "Line 100\nLine 500\nLine 1000"
        assert self.get_output() == expected
    
    def test_program_with_comments(self):
        """Test program with REM statements (comments)."""
        program = """10 REM This is a comment
20 PRINT "This should print"
30 REM Another comment
40 PRINT "This should also print" """
        
        self.capture_output()
        self.interpreter.run(program)
        self.restore_output()
        
        expected = "This should print\nThis should also print"
        assert self.get_output() == expected
    
    def test_end_statement(self):
        """Test END statement terminates program execution."""
        program = """10 PRINT "Before END"
20 END
30 PRINT "After END - should not print" """
        
        self.capture_output()
        self.interpreter.run(program)
        self.restore_output()
        
        assert self.get_output() == "Before END"
    
    def test_error_handling_invalid_line_number(self):
        """Test error handling for invalid GOTO target."""
        program = """10 PRINT "Start"
20 GOTO 999
30 PRINT "End" """
        
        with pytest.raises(Exception) as exc_info:
            self.interpreter.run(program)
        
        assert "line number" in str(exc_info.value).lower()
    
    def test_error_handling_syntax_error(self):
        """Test error handling for syntax errors."""
        program = """10 PRINT "Valid line"
20 INVALID_COMMAND
30 PRINT "Another valid line" """
        
        with pytest.raises(Exception) as exc_info:
            self.interpreter.run(program)
        
        assert "syntax" in str(exc_info.value).lower() or "unknown" in str(exc_info.value).lower()
    
    def test_complex_program_factorial(self):
        """Test a complex program that calculates factorial."""
        program = """10 INPUT "Enter a number: "; N
20 LET F = 1
30 FOR I = 1 TO N
40 LET F = F * I
50 NEXT I
60 PRINT "Factorial of "; N; " is "; F
70 END"""
        
        # Mock input
        original_input = input
        inputs = iter(["5"])
        input_func = lambda prompt="": next(inputs)
        
        try:
            import builtins
            builtins.input = input_func
            
            self.capture_output()
            self.interpreter.run(program)
            self.restore_output()
            
            output = self.get_output()
            assert "120" in output  # 5! = 120
        finally:
            builtins.input = original_input
    
    def test_program_state_isolation(self):
        """Test that each program run starts with clean state."""
        program1 = "10 LET A = 42"
        program2 = "10 PRINT A"
        
        # Run first program
        self.interpreter.run(program1)
        
        # Run second program - should fail if state persists
        with pytest.raises(Exception):
            self.interpreter.run(program2)
    
    def test_load_program_from_string(self):
        """Test loading program from multi-line string."""
        program = """10 PRINT "Line 1"
20 PRINT "Line 2"
30 PRINT "Line 3" """
        
        self.interpreter.load_program(program)
        
        self.capture_output()
        self.interpreter.execute()
        self.restore_output()
        
        expected = "Line 1\nLine 2\nLine 3"
        assert self.get_output() == expected