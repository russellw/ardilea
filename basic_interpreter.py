"""Minimal BASIC interpreter stub for testing purposes."""

class BasicInterpreter:
    """A BASIC interpreter with line number support."""
    
    def __init__(self):
        """Initialize the interpreter."""
        self.program = {}
        self.variables = {}
        self.program_counter = 0
        self.line_numbers = []
        self.for_stack = []  # Stack for FOR loops
        self.return_stack = []  # Stack for GOSUB/RETURN
    
    def run(self, program_text):
        """Run a BASIC program from text."""
        self.load_program(program_text)
        self.execute()
    
    def load_program(self, program_text):
        """Load a BASIC program from text."""
        self.program = {}
        self.variables = {}
        self.for_stack = []
        self.return_stack = []
        
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
            raise Exception(f"Syntax error: Unknown command '{statement}'")
        
        return True
    
    def _execute_print(self, statement):
        """Execute PRINT statement."""
        expr = statement[5:].strip()
        
        if not expr:
            print()  # Empty PRINT statement
            return
        
        # Handle semicolon-separated expressions
        parts = self._parse_print_parts(expr)
        output_parts = []
        
        for i, part in enumerate(parts):
            part = part.strip()
            if part == ';':
                continue
            
            if part.startswith('"') and part.endswith('"'):
                # String literal
                output_parts.append(part[1:-1])
            else:
                # Expression or variable
                try:
                    result = self._evaluate_expression(part)
                    if isinstance(result, float) and result.is_integer():
                        result = int(result)
                    output_parts.append(str(result))
                except Exception as e:
                    raise Exception(f"Error evaluating expression '{part}': {e}")
        
        # Join with spaces (BASIC behavior for semicolon separation)
        print(' '.join(output_parts))
    
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
                raise Exception(f"Undefined line number {target_line} in GOTO statement")
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
        """Execute FOR statement."""
        # Parse: FOR VAR = START TO END [STEP STEP_VALUE]
        parts = statement[3:].strip().split()
        if len(parts) < 5 or parts[2] != '=' or parts[4] != 'TO':
            raise Exception("Invalid FOR syntax")
        
        var_name = parts[1]
        start_value = self._evaluate_expression(parts[3])
        end_value = self._evaluate_expression(parts[5])
        
        step_value = 1
        if len(parts) >= 8 and parts[6] == 'STEP':
            step_value = self._evaluate_expression(parts[7])
        
        # Set loop variable
        self.variables[var_name] = start_value
        
        # Push loop info onto stack
        current_line = self.line_numbers[self.program_counter]
        self.for_stack.append({
            'var': var_name,
            'end': end_value,
            'step': step_value,
            'line': current_line
        })
    
    def _execute_next(self, statement):
        """Execute NEXT statement."""
        if not self.for_stack:
            raise Exception("NEXT without FOR")
        
        # Get loop variable name
        var_name = None
        if len(statement) > 4:
            var_name = statement[4:].strip()
        
        loop_info = self.for_stack[-1]
        
        if var_name and var_name != loop_info['var']:
            raise Exception(f"NEXT {var_name} doesn't match FOR {loop_info['var']}")
        
        # Increment loop variable
        current_value = self.variables[loop_info['var']]
        new_value = current_value + loop_info['step']
        self.variables[loop_info['var']] = new_value
        
        # Check if loop should continue
        if (loop_info['step'] > 0 and new_value <= loop_info['end']) or \
           (loop_info['step'] < 0 and new_value >= loop_info['end']):
            # Continue loop - jump back to FOR line
            for_line = loop_info['line']
            self.program_counter = self.line_numbers.index(for_line)
        else:
            # End loop
            self.for_stack.pop()
    
    def _execute_input(self, statement):
        """Execute INPUT statement."""
        expr = statement[5:].strip()
        
        # Handle different INPUT formats
        if ';' in expr:
            # Format: INPUT "prompt"; variable
            parts = expr.split(';', 1)
            prompt = parts[0].strip()
            var_name = parts[1].strip()
            
            if prompt.startswith('"') and prompt.endswith('"'):
                prompt = prompt[1:-1]
                value = input(prompt)
            else:
                value = input()
        else:
            # Format: INPUT variable
            var_name = expr
            value = input("? ")
        
        # Try to convert to number, otherwise keep as string
        try:
            if '.' in value:
                self.variables[var_name] = float(value)
            else:
                self.variables[var_name] = int(value)
        except ValueError:
            self.variables[var_name] = value
    
    def _evaluate_expression(self, expr):
        """Evaluate a mathematical or string expression."""
        expr = expr.strip()
        
        # String literal
        if expr.startswith('"') and expr.endswith('"'):
            return expr[1:-1]
        
        # Variable (including string variables ending with $)
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
        
        # Handle arithmetic with proper operator precedence
        return self._evaluate_arithmetic(expr)
    
    def _evaluate_arithmetic(self, expr):
        """Evaluate arithmetic expression with proper precedence."""
        expr = expr.strip()
        
        # Handle addition and subtraction (lowest precedence)
        for i in range(len(expr) - 1, -1, -1):
            if expr[i] in ['+', '-'] and i > 0:
                # Make sure it's not a unary operator
                if expr[i-1] not in ['*', '/', '+', '-', '(', '=', '<', '>']:
                    left = self._evaluate_expression(expr[:i])
                    right = self._evaluate_expression(expr[i+1:])
                    if expr[i] == '+':
                        return left + right
                    else:
                        return left - right
        
        # Handle multiplication and division (higher precedence)
        for i in range(len(expr) - 1, -1, -1):
            if expr[i] in ['*', '/']:
                left = self._evaluate_expression(expr[:i])
                right = self._evaluate_expression(expr[i+1:])
                if expr[i] == '*':
                    return left * right
                else:
                    if right == 0:
                        raise Exception("Division by zero")
                    return left / right
        
        # Handle parentheses
        if '(' in expr and ')' in expr:
            start = expr.find('(')
            count = 1
            end = start + 1
            while end < len(expr) and count > 0:
                if expr[end] == '(':
                    count += 1
                elif expr[end] == ')':
                    count -= 1
                end += 1
            
            if count == 0:
                inner = expr[start+1:end-1]
                result = self._evaluate_expression(inner)
                new_expr = expr[:start] + str(result) + expr[end:]
                if new_expr != expr:
                    return self._evaluate_expression(new_expr)
        
        # Single variable or number
        if expr in self.variables:
            return self.variables[expr]
        
        try:
            if '.' in expr:
                return float(expr)
            else:
                return int(expr)
        except ValueError:
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
    
    def _parse_print_parts(self, expr):
        """Parse PRINT statement parts separated by semicolons."""
        parts = []
        current_part = ""
        in_quotes = False
        
        i = 0
        while i < len(expr):
            char = expr[i]
            
            if char == '"':
                in_quotes = not in_quotes
                current_part += char
            elif char == ';' and not in_quotes:
                if current_part.strip():
                    parts.append(current_part.strip())
                current_part = ""
            else:
                current_part += char
            
            i += 1
        
        if current_part.strip():
            parts.append(current_part.strip())
        
        return parts