from clang.cindex import Index, CursorKind, TranslationUnit
from collections import defaultdict
import os

def extract_call_graph(node, current_function=None, edges=None):
    if edges is None:
        edges = []

    if node.kind == CursorKind.FUNCTION_DECL:
        current_function = node.spelling
    elif node.kind == CursorKind.CALL_EXPR and current_function:
        callee = node.spelling
        if callee:
            edges.append((current_function, callee))

    for child in node.get_children():
        extract_call_graph(child, current_function, edges)

    return edges

def build_adj_list(edges):
    adj = defaultdict(list)
    for u, v in edges:
        adj[u].append(v)
    return adj

def has_path_without_edge(adj, start, target, disabled_edge):
    """DFS to check if 'target' is reachable from 'start' without using 'disabled_edge'"""
    visited = set()
    stack = [start]

    while stack:
        curr = stack.pop()
        if curr == target and curr != start:
            return True

        if curr not in visited:
            visited.add(curr)
            for neighbor in adj[curr]:
                
                if (curr, neighbor) == disabled_edge:
                    continue
                if neighbor not in visited:
                    stack.append(neighbor)

    return False

def transitive_reduction(edges):
    adj = build_adj_list(edges)
    reduced_edges = []

    for u, v in edges:
        # Check if an alternative path exists from u to v without direct edge (u, v)
        if not has_path_without_edge(adj, u, v, disabled_edge=(u, v)):
            reduced_edges.append((u, v))

    return reduced_edges

def export_to_dot(edges, filename="graph_reduced.dot"):
    with open(filename, "w") as f:
        f.write("digraph CallGraph {\n")
        f.write("  node [shape=box, fontname=\"Helvetica\", style=filled, fillcolor=\"#f4f4f4\"];\n")
        for caller, callee in edges:
            f.write(f'  "{caller}" -> "{callee}";\n')
        f.write("}\n")
    print(f"\n[SUCCESS] Reduced Graphviz file generated: {filename}")


def extract_includes(node, current_file, edges):
    """
    Traverses the AST to find #include directives within the target project files.
    """
    if node.kind == CursorKind.INCLUSION_DIRECTIVE:
        included_file = node.spelling
        # External system headers (<iostream>, <vector>) filter out 
        # only local project headers (".h", ".hpp") tracked
        if not included_file.startswith("<") and node.location.file:
            source_file = os.path.basename(node.location.file.name)
            target_file = os.path.basename(included_file)

            print("f    Match found: {source_file} includes {target_file}")
            edges.append((source_file, target_file))

    for child in node.get_children():
        extract_includes(child, current_file, edges)


def scan_directory(directory_path):
    index = Index.create()
    all_edges = []

    
    for root, _, files in os.walk(directory_path):
        for file in files:
            if file.endswith(('.cpp', '.c', '.hpp', '.h')):
                full_path = os.path.join(root, file)
                print(f"[PARSING FILE]: {file}")
                
                
                tu = index.parse(full_path, options=TranslationUnit.PARSE_DETAILED_PROCESSING_RECORD)
                extract_includes(tu.cursor, file, all_edges)

    return all_edges

def export_file_graph(edges, filename="file_dependencies.dot"):
    with open(filename, "w") as f:
        f.write("digraph FileDependencies {\n")
        f.write("  node [shape=ellipse, fontname=\"Helvetica\", style=filled, fillcolor=\"#E3F2FD\"];\n")
        for source, target in edges:
            f.write(f'  "{source}" -> "{target}";\n')
        f.write("}\n")

def main():
    index = Index.create()
    translation_unit = index.parse('main.cpp')
    
    raw_edges_func = extract_call_graph(translation_unit.cursor)
    unique_edges_func = list(dict.fromkeys(raw_edges_func))

    print("--- ORIGINAL EDGES ---")
    for u, v in unique_edges_func:
        print(f"  {u} --> {v}")

    # Apply Transitive Reduction
    reduced_edges = transitive_reduction(unique_edges_func)

    print("\n--- TRANSITIVE REDUCED EDGES ---")
    for u, v in reduced_edges:
        print(f"  {u} --> {v}")

    # Export clean DAG to DOT
    export_to_dot(reduced_edges, "graph_reduced.dot")


    project_dir = "."  # Current directory scanning
    raw_edges = scan_directory(project_dir)

    unique_edges = list(dict.fromkeys(raw_edges))

    
    print(f" edges : {raw_edges}")
    print(f" edges : {unique_edges}")
    

    print("\n--- EXTRACTED FILE DEPENDENCIES ---")
    for source, target in unique_edges:
        print(f"  {source} --(#include)--> {target}")

    export_file_graph(unique_edges)

if __name__ == "__main__":
    main()