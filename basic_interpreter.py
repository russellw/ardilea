"""Minimal BASIC interpreter stub for testing purposes."""

class BasicInterpreter:
    """A stub implementation of a BASIC interpreter with line number support."""
    
    def __init__(self):
        """Initialize the interpreter."""
        self.program = {}
        self.variables = {}
        self.program_counter = 0
        self.line_numbers = []
    
    def run(self, program_text):
        """Run a BASIC program from text."""
        self.load_program(program_text)
        self.execute()
    
    def load_program(self, program_text):
        """Load a BASIC program from text."""
        self.program = {}
        self.variables = {}
        
        for line in program_text.strip().split('\n'):
            line = line.strip()
            if not line:
                continue
                
            # Parse line number and statement
            parts = line.split(' ', 1)
            if len(parts) < 2:
                continue
                
            try:
                line_num = int(parts[0])
                statement = parts[1]
                self.program[line_num] = statement
            except ValueError:
                # Skip lines that don't start with a number
                continue
        
        self.line_numbers = sorted(self.program.keys())
    
    def execute(self):
        """Execute the loaded program."""
        if not self.program:
            return
        
        self.program_counter = 0
        
        while self.program_counter < len(self.line_numbers):
            line_num = self.line_numbers[self.program_counter]
            statement = self.program[line_num]
            
            try:
                if not self._execute_statement(statement):
                    break  # END statement or error
                self.program_counter += 1
            except Exception as e:
                raise Exception(f"Error at line {line_num}: {e}")
    
    def _execute_statement(self, statement):
        """Execute a single BASIC statement. Returns False to stop execution."""
        statement = statement.strip()
        
        if statement.startswith('PRINT'):
            self._execute_print(statement)
        elif statement.startswith('LET'):
            self._execute_let(statement)
        elif statement.startswith('GOTO'):
            self._execute_goto(statement)
        elif statement.startswith('IF'):
            self._execute_if(statement)
        elif statement.startswith('FOR'):
            self._execute_for(statement)
        elif statement.startswith('NEXT'):
            self._execute_next(statement)
        elif statement.startswith('INPUT'):
            self._execute_input(statement)
        elif statement.startswith('REM'):
            pass  # Comment, do nothing
        elif statement.startswith('END'):
            return False
        else:
            raise Exception(f"Unknown command: {statement}")
        
        return True
    
    def _execute_print(self, statement):
        """Execute PRINT statement."""
        # Simple implementation - just print everything after PRINT
        expr = statement[5:].strip()
        
        if expr.startswith('"') and expr.endswith('"'):
            # String literal
            print(expr[1:-1])
        else:
            # Expression or variable
            try:
                result = self._evaluate_expression(expr)
                print(result)
            except:
                # If evaluation fails, treat as string
                print(expr)
    
    def _execute_let(self, statement):
        """Execute LET statement."""
        # Simple implementation: LET VAR = VALUE
        parts = statement[3:].strip().split('=', 1)
        if len(parts) != 2:
            raise Exception("Invalid LET syntax")
        
        var_name = parts[0].strip()
        value_expr = parts[1].strip()
        
        value = self._evaluate_expression(value_expr)
        self.variables[var_name] = value
    
    def _execute_goto(self, statement):
        """Execute GOTO statement."""
        line_num_str = statement[4:].strip()
        try:
            target_line = int(line_num_str)
            if target_line in self.line_numbers:
                self.program_counter = self.line_numbers.index(target_line) - 1
            else:
                raise Exception(f"Line number {target_line} not found")
        except ValueError:
            raise Exception("Invalid GOTO syntax")
    
    def _execute_if(self, statement):
        """Execute IF statement."""
        # Simplified IF implementation
        parts = statement[2:].strip().split(' THEN ', 1)
        if len(parts) != 2:
            raise Exception("Invalid IF syntax")
        
        condition = parts[0].strip()
        then_part = parts[1].strip()
        
        if self._evaluate_condition(condition):
            self._execute_statement(then_part)
    
    def _execute_for(self, statement):
        """Execute FOR statement (stub)."""
        # Simplified implementation
        pass
    
    def _execute_next(self, statement):
        """Execute NEXT statement (stub)."""
        # Simplified implementation
        pass
    
    def _execute_input(self, statement):
        """Execute INPUT statement."""
        # Simplified implementation
        parts = statement[5:].strip().split(';')
        if len(parts) == 2:
            prompt = parts[0].strip()
            var_name = parts[1].strip()
            if prompt.startswith('"') and prompt.endswith('"'):
                prompt = prompt[1:-1]
            value = input(prompt)
            try:
                self.variables[var_name] = float(value)
            except ValueError:
                self.variables[var_name] = value
    
    def _evaluate_expression(self, expr):
        """Evaluate a simple expression."""
        expr = expr.strip()
        
        # String literal
        if expr.startswith('"') and expr.endswith('"'):
            return expr[1:-1]
        
        # Variable
        if expr in self.variables:
            return self.variables[expr]
        
        # Number
        try:
            if '.' in expr:
                return float(expr)
            else:
                return int(expr)
        except ValueError:
            pass
        
        # Simple arithmetic (very basic)
        for op in ['+', '-', '*', '/']:
            if op in expr:
                parts = expr.split(op, 1)
                if len(parts) == 2:
                    left = self._evaluate_expression(parts[0])
                    right = self._evaluate_expression(parts[1])
                    if op == '+':
                        return left + right
                    elif op == '-':
                        return left - right
                    elif op == '*':
                        return left * right
                    elif op == '/':
                        return left / right
        
        raise Exception(f"Cannot evaluate expression: {expr}")
    
    def _evaluate_condition(self, condition):
        """Evaluate a simple condition."""
        condition = condition.strip()
        
        for op in ['>', '<', '=']:
            if op in condition:
                parts = condition.split(op, 1)
                if len(parts) == 2:
                    left = self._evaluate_expression(parts[0])
                    right = self._evaluate_expression(parts[1])
                    if op == '>':
                        return left > right
                    elif op == '<':
                        return left < right
                    elif op == '=':
                        return left == right
        
        return False