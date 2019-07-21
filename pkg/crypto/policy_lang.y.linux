%{
#include <stdio.h>
#include <ctype.h>
#include <string.h>
#include <stdlib.h>
#include <stdint.h>
#include <glib.h>
#include <pbc.h>

#include "common.h"
#include "policy_lang.h"

typedef struct
{
	uint64_t value;
	int bits; /* zero if this is a flexint */
}
sized_integer_t;

typedef struct
{
	int k;               /* one if leaf, otherwise threshold */
	char* attr;          /* attribute string if leaf, otherwise null */
	GPtrArray* children; /* pointers to bswabe_policy_t's, len == 0 for leaves */
}
cpabe_policy_t;

cpabe_policy_t* final_policy = 0;

int yylex();
void yyerror( const char* s );
sized_integer_t* expint( uint64_t value, uint64_t bits );
sized_integer_t* flexint( uint64_t value );
cpabe_policy_t* leaf_policy( char* attr );
cpabe_policy_t* kof2_policy( int k, cpabe_policy_t* l, cpabe_policy_t* r );
cpabe_policy_t* kof_policy( int k, GPtrArray* list );
cpabe_policy_t* eq_policy( sized_integer_t* n, char* attr );
cpabe_policy_t* lt_policy( sized_integer_t* n, char* attr );
cpabe_policy_t* gt_policy( sized_integer_t* n, char* attr );
cpabe_policy_t* le_policy( sized_integer_t* n, char* attr );
cpabe_policy_t* ge_policy( sized_integer_t* n, char* attr );
%}

%union
{
	char* str;
	uint64_t nat;
  sized_integer_t* sint;
	cpabe_policy_t* tree;
	GPtrArray* list;
}

%token <str>  TAG
%token <nat>  INTLIT
%type  <sint> number
%type  <tree> policy
%type  <list> arg_list

%left OR
%left AND
%token OF
%token LEQ
%token GEQ

%%

result: policy { final_policy = $1; }

number:   INTLIT '#' INTLIT          { $$ = expint($1, $3); }
        | INTLIT                     { $$ = flexint($1);    }

policy:   TAG                        { $$ = leaf_policy($1);        }
        | policy OR  policy          { $$ = kof2_policy(1, $1, $3); }
        | policy AND policy          { $$ = kof2_policy(2, $1, $3); }
        | INTLIT OF '(' arg_list ')' { $$ = kof_policy($1, $4);     }
        | TAG '=' number             { $$ = eq_policy($3, $1);      }
        | TAG '<' number             { $$ = lt_policy($3, $1);      }
        | TAG '>' number             { $$ = gt_policy($3, $1);      }
        | TAG LEQ number             { $$ = le_policy($3, $1);      }
        | TAG GEQ number             { $$ = ge_policy($3, $1);      }
        | number '=' TAG             { $$ = eq_policy($1, $3);      }
        | number '<' TAG             { $$ = gt_policy($1, $3);      }
        | number '>' TAG             { $$ = lt_policy($1, $3);      }
        | number LEQ TAG             { $$ = ge_policy($1, $3);      }
        | number GEQ TAG             { $$ = le_policy($1, $3);      }
        | '(' policy ')'             { $$ = $2;                     }

arg_list: policy                     { $$ = g_ptr_array_new();
                                       g_ptr_array_add($$, $1); }
        | arg_list ',' policy        { $$ = $1;
                                       g_ptr_array_add($$, $3); }
;

%%

sized_integer_t*
expint( uint64_t value, uint64_t bits )
{
	sized_integer_t* s;

	if( bits == 0 )
		die("error parsing policy: zero-length integer \"%llub%llu\"\n",
				value, bits);
	else if( bits > 64 )
		die("error parsing policy: no more than 64 bits allowed \"%llub%llu\"\n",
				value, bits);

	s = malloc(sizeof(sized_integer_t));
	s->value = value;
	s->bits = bits;

	return s;
}

sized_integer_t*
flexint( uint64_t value )
{
	sized_integer_t* s;

	s = malloc(sizeof(sized_integer_t));
	s->value = value;
	s->bits = 0;

	return s;
}

void
policy_free( cpabe_policy_t* p )
{
	int i;

	if( p->attr )
		free(p->attr);

	for( i = 0; i < p->children->len; i++ )
		policy_free(g_ptr_array_index(p->children, i));
	g_ptr_array_free(p->children, 1);

	free(p);
}

cpabe_policy_t*
leaf_policy( char* attr )
{
	cpabe_policy_t* p;

	p = (cpabe_policy_t*) malloc(sizeof(cpabe_policy_t));
	p->k = 1;
	p->attr = attr;
	p->children = g_ptr_array_new();

	return p;
}

cpabe_policy_t*
kof2_policy( int k, cpabe_policy_t* l, cpabe_policy_t* r )
{
	cpabe_policy_t* p;

	p = (cpabe_policy_t*) malloc(sizeof(cpabe_policy_t));
	p->k = k;
	p->attr = 0;
	p->children = g_ptr_array_new();
	g_ptr_array_add(p->children, l);
	g_ptr_array_add(p->children, r);

	return p;
}

cpabe_policy_t*
kof_policy( int k, GPtrArray* list )
{
	cpabe_policy_t* p;

	if( k < 1 )
		die("error parsing policy: trivially satisfied operator \"%dof\"\n", k);
	else if( k > list->len )
		die("error parsing policy: unsatisfiable operator \"%dof\" (only %d operands)\n",
				k, list->len);
	else if( list->len == 1 )
		die("error parsing policy: identity operator \"%dof\" (only one operand)\n", k);

	p = (cpabe_policy_t*) malloc(sizeof(cpabe_policy_t));
	p->k = k;
	p->attr = 0;
	p->children = list;

	return p;
}

char*
bit_marker( char* base, char* tplate, int bit, char val )
{
	char* lx;
	char* rx;
	char* s;

 	lx = g_strnfill(64 - bit - 1, 'x');
	rx = g_strnfill(bit, 'x');
	s = g_strdup_printf(tplate, base, lx, !!val, rx);
	free(lx);
	free(rx);

	return s;
}

cpabe_policy_t*
eq_policy( sized_integer_t* n, char* attr )
{
	if( n->bits == 0 )
		return leaf_policy
			(g_strdup_printf("%s_flexint_%llu", attr, n->value));
	else
		return leaf_policy
			(g_strdup_printf("%s_expint%02d_%llu", attr, n->bits, n->value));

	return 0;
}

cpabe_policy_t*
bit_marker_list( int gt, char* attr, char* tplate, int bits, uint64_t value )
{
	cpabe_policy_t* p;
	int i;

	i = 0;
	while( gt ? (((uint64_t)1)<<i & value) : !(((uint64_t)1)<<i & value) )
		i++;

	p = leaf_policy(bit_marker(attr, tplate, i, gt));
	for( i = i + 1; i < bits; i++ )
		if( gt )
			p = kof2_policy(((uint64_t)1<<i & value) ? 2 : 1, p,
											leaf_policy(bit_marker(attr, tplate, i, gt)));
		else
			p = kof2_policy(((uint64_t)1<<i & value) ? 1 : 2, p,
											leaf_policy(bit_marker(attr, tplate, i, gt)));

	return p;
}

cpabe_policy_t*
flexint_leader( int gt, char* attr, uint64_t value )
{
	cpabe_policy_t* p;
	int k;

	p = (cpabe_policy_t*) malloc(sizeof(cpabe_policy_t));
	p->attr = 0;
	p->children = g_ptr_array_new();

	for( k = 2; k <= 32; k *= 2 )
		if( ( gt && ((uint64_t)1<<k) >  value) ||
				(!gt && ((uint64_t)1<<k) >= value) )
			g_ptr_array_add
				(p->children, leaf_policy
				 (g_strdup_printf(gt ? "%s_ge_2^%02d" : "%s_lt_2^%02d", attr, k)));

	p->k = gt ? 1 : p->children->len;

	if( p->children->len == 0 )
	{
		policy_free(p);
		p = 0;
	}
	else if( p->children->len == 1 )
	{
		cpabe_policy_t* t;
		
		t = g_ptr_array_remove_index(p->children, 0);
		policy_free(p);
		p = t;
	}
	
	return p;
}

cpabe_policy_t*
cmp_policy( sized_integer_t* n, int gt, char* attr )
{
	cpabe_policy_t* p;
	char* tplate;

	/* some error checking */

	if( gt && n->value >= ((uint64_t)1<<(n->bits ? n->bits : 64)) - 1 )
		die("error parsing policy: unsatisfiable integer comparison %s > %llu\n"
				"(%d-bits are insufficient to satisfy)\n", attr, n->value,
				n->bits ? n->bits : 64);
	else if( !gt && n->value == 0 )
		die("error parsing policy: unsatisfiable integer comparison %s < 0\n"
				"(all numerical attributes are unsigned)\n", attr);
	else if( !gt && n->value > ((uint64_t)1<<(n->bits ? n->bits : 64)) - 1 )
		die("error parsing policy: trivially satisfied integer comparison %s < %llu\n"
				"(any %d-bit number will satisfy)\n", attr, n->value,
				n->bits ? n->bits : 64);

	/* create it */

	/* horrible */
	tplate = n->bits ?
		g_strdup_printf("%%s_expint%02d_%%s%%d%%s", n->bits) :
		strdup("%s_flexint_%s%d%s");
	p = bit_marker_list(gt, attr, tplate, n->bits ? n->bits :
											(n->value >= ((uint64_t)1<<32) ? 64 :
											 n->value >= ((uint64_t)1<<16) ? 32 :
											 n->value >= ((uint64_t)1<< 8) ? 16 :
											 n->value >= ((uint64_t)1<< 4) ?  8 :
											 n->value >= ((uint64_t)1<< 2) ?  4 : 2), n->value);
	free(tplate);

	if( !n->bits )
	{
		cpabe_policy_t* l;
		
		l = flexint_leader(gt, attr, n->value);
		if( l )
			p = kof2_policy(gt ? 1 : 2, l, p);
	}

	return p;
}

cpabe_policy_t*
lt_policy( sized_integer_t* n, char* attr )
{
	return cmp_policy(n, 0, attr);
}

cpabe_policy_t*
gt_policy( sized_integer_t* n, char* attr )
{
	return cmp_policy(n, 1, attr);
}

cpabe_policy_t*
le_policy( sized_integer_t* n, char* attr )
{
	n->value++;
	return cmp_policy(n, 0, attr);
}

cpabe_policy_t*
ge_policy( sized_integer_t* n, char* attr )
{
	n->value--;
	return cmp_policy(n, 1, attr);
}

char* cur_string = 0;

#define PEEK_CHAR ( *cur_string ? *cur_string     : EOF )
#define NEXT_CHAR ( *cur_string ? *(cur_string++) : EOF )

int
yylex()
{
  int c;
	int r;

  while( isspace(c = NEXT_CHAR) );

	r = 0;
  if( c == EOF )
    r = 0;
	else if( c == '&' )
		r = AND;
	else if( c == '|' )
		r = OR;
	else if( strchr("(),=#", c) || (strchr("<>", c) && PEEK_CHAR != '=') )
		r = c;
	else if( c == '<' && PEEK_CHAR == '=' )
	{
		NEXT_CHAR;
		r = LEQ;
	}
	else if( c == '>' && PEEK_CHAR == '=' )
	{
		NEXT_CHAR;
		r = GEQ;
	}
	else if( isdigit(c) )
	{
		GString* s;

		s = g_string_new("");
		g_string_append_c(s, c);
		while( isdigit(PEEK_CHAR) )
			g_string_append_c(s, NEXT_CHAR);

		sscanf(s->str, "%llu", &(yylval.nat));

		g_string_free(s, 1);
		r = INTLIT;
	}
	else if( isalpha(c) )
	{
		GString* s;

		s = g_string_new("");
		g_string_append_c(s, c);

		while( isalnum(PEEK_CHAR) || PEEK_CHAR == '_' )
			g_string_append_c(s, NEXT_CHAR);

		if( !strcmp(s->str, "and") )
		{
			g_string_free(s, 1);
			r = AND;
		}
		else if( !strcmp(s->str, "or") )
		{
			g_string_free(s, 1);
			r = OR;
		}
		else if( !strcmp(s->str, "of") )
		{
			g_string_free(s, 1);
			r = OF;
		}
		else
		{
			yylval.str = s->str;
			g_string_free(s, 0);
			r = TAG;
		}
	}
	else
		die("syntax error at \"%c%s\"\n", c, cur_string);

	return r;
}

void
yyerror( const char* s )
{
  die("error parsing policy: %s\n", s);
}

#define POLICY_IS_OR(p)  (((cpabe_policy_t*)(p))->k == 1 && ((cpabe_policy_t*)(p))->children->len)
#define POLICY_IS_AND(p) (((cpabe_policy_t*)(p))->k == ((cpabe_policy_t*)(p))->children->len)

void
merge_child( cpabe_policy_t* p, int i )
{
	int j;
	cpabe_policy_t* c;

	c = g_ptr_array_index(p->children, i);
	if( POLICY_IS_AND(p) )
	{
		p->k += c->k;
		p->k--;
	}

	g_ptr_array_remove_index_fast(p->children, i);
	for( j = 0; j < c->children->len; j++ )
		g_ptr_array_add(p->children, g_ptr_array_index(c->children, j));

	g_ptr_array_free(c->children, 0);
	free(c);
}

void
simplify( cpabe_policy_t* p )
{
	int i;

	for( i = 0; i < p->children->len; i++ )
		simplify(g_ptr_array_index(p->children, i));

	if( POLICY_IS_OR(p) )
		for( i = 0; i < p->children->len; i++ )
			if( POLICY_IS_OR(g_ptr_array_index(p->children, i)) )
				merge_child(p, i);

	if( POLICY_IS_AND(p) )
		for( i = 0; i < p->children->len; i++ )
			if( POLICY_IS_AND(g_ptr_array_index(p->children, i)) )
				merge_child(p, i);
}

int
cmp_tidy( const void* a, const void* b )
{
	cpabe_policy_t* pa;
	cpabe_policy_t* pb;

	pa = *((cpabe_policy_t**) a);
	pb = *((cpabe_policy_t**) b);

	if(      pa->children->len >  0 && pb->children->len == 0 )
		return -1;
	else if( pa->children->len == 0 && pb->children->len >  0 )
		return 1;
	else if( pa->children->len == 0 && pb->children->len == 0 )
		return strcmp(pa->attr, pb->attr);
	else
		return 0;	
}

void
tidy( cpabe_policy_t* p )
{
	int i;

	for( i = 0; i < p->children->len; i++ )
		tidy(g_ptr_array_index(p->children, i));

	if( p->children->len > 0 )
		qsort(p->children->pdata, p->children->len,
					sizeof(cpabe_policy_t*), cmp_tidy);
}

char*
format_policy_postfix( cpabe_policy_t* p )
{
	int i;
	char* r;
	char* s;
	char* t;

	if( p->children->len == 0 )
		return strdup(p->attr);

	r = format_policy_postfix(g_ptr_array_index(p->children, 0));
	for( i = 1; i < p->children->len; i++ )
	{
		s = format_policy_postfix(g_ptr_array_index(p->children, i));
		t = g_strjoin(" ", r, s, (char*) 0);
		free(r);
		free(s);
		r = t;
	}
	
	t = g_strdup_printf("%s %dof%d", r, p->k, p->children->len);
 	free(r);

	return t;
}

/*
	Crufty.
*/
int
actual_bits( uint64_t value )
{
	int i;

	for( i = 32; i >= 1; i /= 2 )
		if( value >= ((uint64_t)1<<i) )
			return i * 2;

	return 1;
}

/*
	It is pretty crufty having this here since it is only used in
	keygen. Maybe eventually there will be a separate .c file with the
	policy_lang module.
*/
void
parse_attribute( GSList** l, char* a )
{
	if( !strchr(a, '=') )
		*l = g_slist_append(*l, a);
	else
	{
		int i;
		char* s;
		char* tplate;
		uint64_t value;
		int bits;

		s = malloc(sizeof(char) * strlen(a));

		if( sscanf(a, " %s = %llu # %u ", s, &value, &bits) == 3 )
		{
			/* expint */

			if( bits > 64 )
				die("error parsing attribute \"%s\": 64 bits is the maximum allowed\n",
						a, value, bits);

			if( value >= ((uint64_t)1<<bits) )
				die("error parsing attribute \"%s\": value %llu too big for %d bits\n",
						a, value, bits);

			tplate = g_strdup_printf("%%s_expint%02d_%%s%%d%%s", bits);
			for( i = 0; i < bits; i++ )
				*l = g_slist_append
					(*l, bit_marker(s, tplate, i, !!((uint64_t)1<<i & value)));
			free(tplate);

			*l = g_slist_append
				(*l, g_strdup_printf("%s_expint%02d_%llu", s, bits, value));
		}
		else if( sscanf(a, " %s = %llu ", s, &value) == 2 )
		{
			/* flexint */

			for( i = 2; i <= 32; i *= 2 )
				*l = g_slist_append
					(*l, g_strdup_printf
					 (value < ((uint64_t)1<<i) ? "%s_lt_2^%02d" : "%s_ge_2^%02d", s, i));

			for( i = 0; i < 64; i++ )
				*l = g_slist_append
					(*l, bit_marker(s, "%s_flexint_%s%d%s", i, !!((uint64_t)1<<i & value)));

			*l = g_slist_append
				(*l, g_strdup_printf("%s_flexint_%llu", s, value));
		}
		else
			die("error parsing attribute \"%s\"\n"
					"(note that numerical attributes are unsigned integers)\n",	a);

 		free(s);
	}	
}

char*
parse_policy_lang( char* s )
{
	char* parsed_policy;

	cur_string = s;

	yyparse();
 	simplify(final_policy);
 	tidy(final_policy);
	parsed_policy = format_policy_postfix(final_policy);

	policy_free(final_policy);

	return parsed_policy;
}
